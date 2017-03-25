package shape

import (
	"fmt"
	"os"

	"github.com/momchil-atanasov/go-data-front/decoder/obj"

	"github.com/ironsmile/raytracer/geometry"
)

type Object struct {
	// A model wich contains the parsed .obj information such as objects, meshes, faces and raw
	// vertices.
	model *obj.Model

	// The center of the object
	Center *geometry.Point

	// All the faces which compose this object
	Faces []Shape

	// A sphere which contains all the points of the object faces
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

	prim, distance, normal := IntersectMultiple(o.Faces, ray)
	if prim == nil {
		return MISS, distance, outNormal
	}
	return HIT, distance, normal
}

func (o *Object) GetNormal(pos *geometry.Point) *geometry.Vector {
	//!TODO: implement
	return &geometry.Vector{0, 0, -1}
}

func (o *Object) computeBoundingSphere() error {
	//!TODO: maybe implement one of the following:
	// https://www.inf.ethz.ch/personal/gaertner/texts/own_work/esa99_final.pdf
	// http://www.ep.liu.se/ecp/034/009/ecp083409.pdf
	// A the moment this is a simple and buggy (see next comment) exhaustion search.

	// Bug: this method makes an implicit guess that the object is centered around the o.Center.
	// This might not be true.

	maxRadius := 0.0

	for ind, triangle := range o.Faces {
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
func NewObject(filePath string, center *geometry.Point) (*Object, error) {
	decoderLimits := obj.DefaultLimits()
	decoderLimits.MaxVertexCount = 4294967296
	decoder := obj.NewDecoder(decoderLimits)
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
	o.Faces = make([]Shape, 0, len(model.Vertices)/3+1)
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

				faceShapes, err := o.faceToShapes(model, face)
				if err != nil {
					return nil, fmt.Errorf(
						"face %d [mesh: %d, obj: %s]: %s",
						faceIndex,
						meshIndex,
						obj.Name,
						err,
					)
				}
				for _, faceShape := range faceShapes {
					o.Faces = append(o.Faces, faceShape)
				}
			}
		}
	}

	fmt.Printf("%s has %d triangles\n", filePath, len(o.Faces))

	if err := o.computeBoundingSphere(); err != nil {
		return nil, err
	}

	return o, nil
}

func (o *Object) faceToShapes(model *obj.Model, face *obj.Face) ([]Shape, error) {

	//!TODO: maybe remove this scale factor?
	scaleFactor := 100.0

	switch len(face.References) {
	case 3:
		a := model.Vertices[face.References[0].VertexIndex]
		b := model.Vertices[face.References[1].VertexIndex]
		c := model.Vertices[face.References[2].VertexIndex]

		triangleVertices := [3]geometry.Point{
			*geometry.NewPoint(a.X/scaleFactor, a.Y/scaleFactor, a.Z/scaleFactor).Plus(o.Center),
			*geometry.NewPoint(b.X/scaleFactor, b.Y/scaleFactor, b.Z/scaleFactor).Plus(o.Center),
			*geometry.NewPoint(c.X/scaleFactor, c.Y/scaleFactor, c.Z/scaleFactor).Plus(o.Center),
		}

		return []Shape{NewTriangle(triangleVertices)}, nil
	case 4:
		//!TODO: Load an actual rectangle which would hopefully be faster. Stop
		// loading two triangles.
		// 0, 1, 3
		// 1, 2, 3
		a := model.Vertices[face.References[0].VertexIndex]
		b := model.Vertices[face.References[1].VertexIndex]
		c := model.Vertices[face.References[3].VertexIndex]

		triangleOne := [3]geometry.Point{
			*geometry.NewPoint(a.X/scaleFactor, a.Y/scaleFactor, a.Z/scaleFactor).Plus(o.Center),
			*geometry.NewPoint(b.X/scaleFactor, b.Y/scaleFactor, b.Z/scaleFactor).Plus(o.Center),
			*geometry.NewPoint(c.X/scaleFactor, c.Y/scaleFactor, c.Z/scaleFactor).Plus(o.Center),
		}

		a = model.Vertices[face.References[1].VertexIndex]
		b = model.Vertices[face.References[2].VertexIndex]
		c = model.Vertices[face.References[3].VertexIndex]

		triangleTwo := [3]geometry.Point{
			*geometry.NewPoint(a.X/scaleFactor, a.Y/scaleFactor, a.Z/scaleFactor).Plus(o.Center),
			*geometry.NewPoint(b.X/scaleFactor, b.Y/scaleFactor, b.Z/scaleFactor).Plus(o.Center),
			*geometry.NewPoint(c.X/scaleFactor, c.Y/scaleFactor, c.Z/scaleFactor).Plus(o.Center),
		}

		return []Shape{
			NewTriangle(triangleOne),
			NewTriangle(triangleTwo),
		}, nil
	}

	return nil, fmt.Errorf(
		"don't know how to load %d vertices in a shape",
		len(face.References),
	)
}
