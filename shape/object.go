package shape

import (
	"fmt"
	"os"

	"github.com/momchil-atanasov/go-data-front/decoder/obj"

	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
)

// Object represents a object in 3d space which shape is loaded from a .obj file
type Object struct {
	BasicShape

	// A model wich contains the parsed .obj information such as objects, meshes,
	// faces and raw vertices.
	model *obj.Model

	// All the triangles which compose this object
	Triangles []Shape
}

// Intersect implements the Shape interface
func (o *Object) Intersect(ray geometry.Ray) (int, float64, geometry.Vector) {
	var outNormal geometry.Vector

	prim, distance, normal := IntersectMultiple(o.Triangles, ray)

	if prim == nil {
		return MISS, distance, outNormal
	}

	return HIT, distance, normal
}

// IntersectP implements the Shape interface
func (o *Object) IntersectP(ray geometry.Ray) bool {
	return IntersectPMultiple(o.Triangles, ray)
}

// NewObject parses an .obj file (`filePath`) and returns an Object, which represents
// it. It places the object at the position, given by its second argument - `center`.
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
						"face %d [mesh: %d, obj: %s] has %d points, cannot load it",
						faceIndex, meshIndex, obj.Name, len(face.References))
				}

				a := model.Vertices[face.References[0].VertexIndex]
				b := model.Vertices[face.References[1].VertexIndex]
				c := model.Vertices[face.References[2].VertexIndex]

				triangleVertices := [3]geometry.Vector{
					geometry.NewVector(a.X, a.Y, a.Z),
					geometry.NewVector(b.X, b.Y, b.Z),
					geometry.NewVector(c.X, c.Y, c.Z),
				}

				o.Triangles = append(o.Triangles, NewTriangle(triangleVertices))
			}
		}
	}

	fmt.Printf("%s has %d triangles\n", filePath, len(o.Triangles))

	computedBBox, err := o.objectBound()
	if err != nil {
		return nil, err
	}

	o.bbox = computedBBox

	return o, nil
}

// objectBound calculates a bounding box which encapsulates the shape
func (o *Object) objectBound() (*bbox.BBox, error) {
	var retBox *bbox.BBox
	for ind, obj := range o.Triangles {
		obj, ok := obj.(*Triangle)
		if !ok {
			return nil, fmt.Errorf(
				"a shape in object.Triangles is not a triangle? Index %d", ind)
		}
		if retBox == nil {
			retBox = bbox.FromPoint(obj.Vertices[0])
		}
		for i := 0; i < 3; i++ {
			retBox = bbox.UnionPoint(retBox, obj.Vertices[i])
		}
	}
	if retBox == nil {
		return nil, fmt.Errorf("obj does not have any faces")
	}
	return retBox, nil
}
