package scene

import (
	"fmt"
	"os"

	"github.com/momchil-atanasov/go-data-front/decoder/obj"

	"github.com/ironsmile/raytracer/geometry"
)

type Object struct {
	BasePrimitive

	// Identification string for thins object
	id string

	// A model wich contains the parsed .obj information such as objects, meshes, faces and raw
	// vertices.
	model *obj.Model

	// The center of the object
	Center *geometry.Point

	// All the triangles which compose this object
	Triangles []Primitive

	// A sphere which contains all the points of the object triangles
	boundingSphere *Sphere
}

func (o *Object) GetType() int {
	return OBJECT
}

func (o *Object) Intersect(ray *geometry.Ray, dist float64) (int, float64) {
	if o.boundingSphere != nil {
		spHit, _ := o.boundingSphere.Intersect(ray, dist)
		if spHit == MISS {
			return MISS, dist
		}
	}

	prim, distance := IntersectPrimitives(o.Triangles, ray)
	if prim == nil {
		return MISS, distance
	}
	return HIT, distance
}

func (o *Object) GetNormal(pos *geometry.Point) *geometry.Vector {
	//!TODO: implement
	return &geometry.Vector{0, 0, -1}
}

func (o *Object) String() string {
	return fmt.Sprintf("Object <%s>", o.id)
}

func (o *Object) computeBoundingSphere() error {
	//!TODO: maybe implement one of the following:
	// https://www.inf.ethz.ch/personal/gaertner/texts/own_work/esa99_final.pdf
	// http://www.ep.liu.se/ecp/034/009/ecp083409.pdf

	// Bug: this method makes an implicit guess that the object is centered around the o.Center.
	// This might not be true.

	maxRadius := 0.0

	for _, triangle := range o.Triangles {
		triangle, ok := triangle.(*Triangle)
		if !ok {
			return fmt.Errorf("A triangle was not actually a triangle? Strange!")
		}
		for _, vertice := range triangle.Vertices {
			distance := geometry.Distance(o.Center, vertice)
			if distance > maxRadius {
				maxRadius = distance
			}
		}
	}

	o.boundingSphere = NewSphere(*o.Center, maxRadius)
	return nil
}

func NewObject(filePath string, center *geometry.Point) (*Object, error) {
	decoder := obj.NewDecoder(obj.DefaultLimits())
	objFile, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer objFile.Close()

	model, err := decoder.Decode(objFile)

	if err != nil {
		return nil, err
	}

	fmt.Printf("model %s has %d models\n", filePath, len(model.Objects))

	o := new(Object)
	o.Center = center
	o.Triangles = make([]Primitive, 0, len(model.Vertices)/3+1)
	o.id = filePath
	o.model = model

	//!TODO: maybe remove this scale factor?
	scaleFactor := 100.0

	for _, obj := range model.Objects {
		fmt.Printf("object %s has %d meshes\n", obj.Name, len(obj.Meshes))
		for meshIndex, mesh := range obj.Meshes {
			fmt.Printf("mesh %d is from %s and has %d faces\n", meshIndex, mesh.MaterialName,
				len(mesh.Faces))
			for faceIndex, face := range mesh.Faces {
				if len(face.References) != 3 {
					return nil, fmt.Errorf(
						"face %d [mesh: %d, obj: %s] has %d points, don't know how to load it",
						faceIndex, meshIndex, obj.Name, len(face.References))
				}

				a := model.Vertices[face.References[0].VertexIndex]
				b := model.Vertices[face.References[1].VertexIndex]
				c := model.Vertices[face.References[2].VertexIndex]

				triangleVertices := [3]*geometry.Point{
					geometry.NewPoint(a.X/scaleFactor, a.Y/scaleFactor, a.Z/scaleFactor).Plus(center),
					geometry.NewPoint(b.X/scaleFactor, b.Y/scaleFactor, b.Z/scaleFactor).Plus(center),
					geometry.NewPoint(c.X/scaleFactor, c.Y/scaleFactor, c.Z/scaleFactor).Plus(center),
				}

				o.Triangles = append(o.Triangles, NewTriangle(triangleVertices))
			}
		}
	}

	fmt.Printf("%s has %d triangles\n", o, len(o.Triangles))

	if err := o.computeBoundingSphere(); err != nil {
		return nil, err
	}

	return o, nil
}
