package shape

import (
	"fmt"
	"os"

	"github.com/momchil-atanasov/go-data-front/decoder/obj"

	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
)

type Object struct {
	// A model wich contains the parsed .obj information such as objects, meshes, faces and raw
	// vertices.
	model *obj.Model

	// The center of the object
	Center *geometry.Point

	// All the triangles which compose this object
	Triangles []Shape

	// A sphere which contains all the points of the object triangles
	boundingSphere *Sphere
}

func (o *Object) Intersect(ray geometry.Ray, dist float64) (int, float64, geometry.Vector) {
	var outNormal geometry.Vector

	if o.boundingSphere != nil {
		spHit, _, _ := o.boundingSphere.Intersect(ray, dist)
		if spHit == MISS {
			return MISS, dist, outNormal
		}
	}

	prim, distance, normal := IntersectMultiple(o.Triangles, ray)
	if prim == nil {
		return MISS, distance, outNormal
	}
	return HIT, distance, normal
}

func (o *Object) computeBoundingSphere() error {
	//!TODO: maybe implement one of the following:
	// https://www.inf.ethz.ch/personal/gaertner/texts/own_work/esa99_final.pdf
	// http://www.ep.liu.se/ecp/034/009/ecp083409.pdf
	// A the moment this is a simple and buggy (see next comment) exhaustion search.

	// Bug: this method makes an implicit guess that the object is centered around the o.Center.
	// This might not be true.

	maxRadius := 0.0

	for ind, triangle := range o.Triangles {
		triangle, ok := triangle.(*Triangle)
		if !ok {
			fmt.Printf("A shape in object.Triangles is not a triangle? Index %d", ind)
			continue
		}
		for _, vertice := range triangle.Vertices {
			distance := geometry.Distance(o.Center, &vertice)
			if distance > maxRadius {
				maxRadius = distance
			}
		}
	}

	o.boundingSphere = NewSphere(*o.Center, maxRadius)
	return nil
}

// NewObject parses an .obj file (`filePath`) and returns an Object, which represents it. It places
// the object at the position, given by its second argument - `center`.
func NewObject(filePath string) (*Object, error) {
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

	o := &Object{}
	o.Center = geometry.NewPoint(0, 0, 0)
	o.Triangles = make([]Shape, 0, len(model.Vertices)/3+1)
	o.model = model

	for _, obj := range model.Objects {
		fmt.Printf("object %s has %d meshes\n", obj.Name, len(obj.Meshes))
		for meshIndex, mesh := range obj.Meshes {
			meshName := mesh.MaterialName
			if len(meshName) < 1 {
				meshName = "Unknown"
			}
			fmt.Printf("mesh %d is from `%s` and has %d faces\n", meshIndex, meshName,
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

				triangleVertices := [3]geometry.Point{
					*geometry.NewPoint(a.X, a.Y, a.Z),
					*geometry.NewPoint(b.X, b.Y, b.Z),
					*geometry.NewPoint(c.X, c.Y, c.Z),
				}

				o.Triangles = append(o.Triangles, NewTriangle(triangleVertices))
			}
		}
	}

	fmt.Printf("%s has %d triangles\n", filePath, len(o.Triangles))

	if err := o.computeBoundingSphere(); err != nil {
		return nil, err
	}

	return o, nil
}

func (o *Object) objectBound() (*bbox.BBox, error) {
	//!todo:
	// create a empty box
	var retBox *bbox.BBox
	for ind, obj := range o.Triangles {
		obj, ok := obj.(*Triangle)
		if !ok {
			return nil, fmt.Errorf("a shape in object.Triangles is not a triangle? Index %d", ind)
		}
		if retBox == nil {
			retBox = retBox.UnionPoint(&obj.Vertices[0])
		}
		for i := 0; i < 3; i++ {
			retBox = retBox.UnionPoint(&obj.Vertices[i])
		}
	}
	if retBox == nil {
		return nil, fmt.Errorf("obj does not have any faces")
	}
	return retBox, nil
}
