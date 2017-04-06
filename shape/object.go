package shape

import (
	"fmt"
	"os"

	"github.com/momchil-atanasov/go-data-front/decoder/obj"

	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
)

type Object struct {
	BasicShape

	// A model wich contains the parsed .obj information such as objects, meshes, faces and raw
	// vertices.
	model *obj.Model

	// The center of the object
	Center *geometry.Point

	// All the triangles which compose this object
	Triangles []Shape
}

func (o *Object) Intersect(ray geometry.Ray, dist float64) (int, float64, geometry.Vector) {
	var outNormal geometry.Vector

	prim, distance, normal := IntersectMultiple(o.Triangles, ray)

	if prim == nil {
		return MISS, dist, outNormal
	}

	if dist < distance {
		return MISS, dist, outNormal
	}

	return HIT, distance, normal
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
			return nil, fmt.Errorf("a shape in object.Triangles is not a triangle? Index %d", ind)
		}
		if retBox == nil {
			retBox = bbox.FromPoint(&obj.Vertices[0])
		}
		for i := 0; i < 3; i++ {
			retBox = bbox.UnionPoint(retBox, &obj.Vertices[i])
		}
	}
	if retBox == nil {
		return nil, fmt.Errorf("obj does not have any faces")
	}
	return retBox, nil
}
