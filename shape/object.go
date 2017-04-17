package shape

import (
	"fmt"
	"os"

	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"

	"github.com/momchil-atanasov/go-data-front/decoder/obj"
)

// Object represents a object in 3d space which shape is loaded from a .obj file
type Object struct {
	BasicShape

	// A model wich contains the parsed .obj information such as objects, meshes,
	// faces and raw vertices.
	model *obj.Model

	// All the meshes which compose this object
	meshes []Shape
}

// Intersect implements the Shape interface
func (o *Object) Intersect(ray geometry.Ray) (int, float64, geometry.Vector) {
	prim, distance, normal := IntersectMultiple(o.meshes, ray)

	if prim == nil {
		return MISS, distance, normal
	}

	return HIT, distance, normal
}

// IntersectP implements the Shape interface
func (o *Object) IntersectP(ray geometry.Ray) bool {
	return IntersectPMultiple(o.meshes, ray)
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
	o.model = model
	var trianglesCount uint32

	for _, modelObj := range model.Objects {
		fmt.Printf("object %s has %d meshes\n", modelObj.Name, len(modelObj.Meshes))
		for meshIndex, mesh := range modelObj.Meshes {
			meshName := mesh.MaterialName
			if len(meshName) < 1 {
				meshName = "Unknown"
			}
			fmt.Printf("mesh %d is from `%s` and has %d faces\n", meshIndex, meshName,
				len(mesh.Faces))
			meshTriangles := make([]Shape, 0, len(mesh.Faces))
			for faceIndex, face := range mesh.Faces {
				if len(face.References) != 3 {
					return nil, fmt.Errorf(
						"face %d [mesh: %d, obj: %s] has %d points, cannot load it",
						faceIndex, meshIndex, modelObj.Name, len(face.References))
				}

				tr := NewMeshTriangle(model, face)
				meshTriangles = append(meshTriangles, tr)
				trianglesCount++
			}

			triagularMesh := NewMesh(model, meshTriangles)
			o.bbox = bbox.Union(o.bbox, triagularMesh.GetObjectBBox())
			o.meshes = append(o.meshes, triagularMesh)
		}
	}

	fmt.Printf("%s has %d triangles\n", filePath, trianglesCount)

	return o, nil
}

func (o *Object) GetAllShapes() []Shape {
	var shapes []Shape
	for _, mesh := range o.meshes {
		shapes = append(shapes, mesh.GetAllShapes()...)
	}
	return shapes
}
