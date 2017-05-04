package shape

import (
	"fmt"
	"os"
	"strings"

	"github.com/ironsmile/raytracer/mat"

	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"

	"github.com/mokiat/go-data-front/decoder/mtl"
	"github.com/mokiat/go-data-front/decoder/obj"
)

const objFileSuffix = ".obj"
const mtlFileSuffix = ".mtl"

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
func (o *Object) Intersect(geometry.Ray, *DifferentialGeometry) bool {
	panic("Object shape is not directly intersectable: Intersect")
}

// IntersectP implements the Shape interface
func (o *Object) IntersectP(geometry.Ray) bool {
	panic("Object shape is not directly intersectable: IntersectP")
}

// NewObject parses an .obj file (`filePath`) and returns an Object, which represents
// it. It places the object at the position, given by its second argument - `center`.
func NewObject(filePath string) (*Object, error) {
	objDecoder := obj.NewDecoder(obj.DefaultLimits())
	mtlDecoder := mtl.NewDecoder(mtl.DefaultLimits())
	objFile, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer objFile.Close()

	model, err := objDecoder.Decode(objFile)

	if err != nil {
		return nil, err
	}

	fmt.Printf("model %s has %d objects\n", filePath, len(model.Objects))

	var matLib *mtl.Library

	if strings.HasSuffix(filePath, objFileSuffix) {
		materialPath := strings.TrimSuffix(filePath, objFileSuffix)
		materialPath += mtlFileSuffix

		if matFile, err := os.Open(materialPath); err == nil {
			defer matFile.Close()
			matLib, err = mtlDecoder.Decode(matFile)

			if err != nil {
				return nil, fmt.Errorf("error decoding material file: %s", err)
			}
		} else {
			fmt.Printf("Error opening material file %s: %s\n", materialPath, err)
		}
	}

	fmt.Printf("model %s has has a material: %p\n", filePath, matLib)

	o := &Object{}
	o.model = model
	var trianglesCount int

	for _, modelObj := range model.Objects {
		fmt.Printf("object %s has %d meshes\n", modelObj.Name, len(modelObj.Meshes))
		for meshIndex, mesh := range modelObj.Meshes {
			meshName := mesh.MaterialName
			if len(meshName) < 1 {
				meshName = "Unknown"
			}
			fmt.Printf("mesh %d is from `%s` and has %d faces\n", meshIndex, meshName,
				len(mesh.Faces))
			trianglesCount += len(mesh.Faces)
			faceMesh := NewMesh(model, mesh)

			if matLib != nil {
				if foundMat, ok := matLib.FindMaterial(mesh.MaterialName); ok {
					faceMesh.SetMaterial(mat.Material{
						Color: geometry.NewColor(
							foundMat.DiffuseColor.R,
							foundMat.DiffuseColor.G,
							foundMat.DiffuseColor.B,
						),
						Diff: foundMat.Dissolve,
						Refr: 1 - foundMat.Dissolve,
					})
				}
			}

			if faceMesh.GetMaterial() == nil {
				faceMesh.SetMaterial(mat.DefaultMetiral())
			}

			o.bbox = bbox.Union(o.bbox, faceMesh.GetObjectBBox())
			o.meshes = append(o.meshes, faceMesh)
		}
	}

	fmt.Printf("%s has %d triangles\n", filePath, trianglesCount)

	return o, nil
}

// CanIntersect implements the Shape interface
func (o *Object) CanIntersect() bool {
	return false
}

// GetMaterial implements the Shape interface
func (o *Object) GetMaterial() *mat.Material {
	panic("GetMaterial should not  be called for shape.Object")
}

// SetMaterial implements Shape interface
func (o *Object) SetMaterial(mat.Material) {
	panic("SetMaterial should not  be called for shape.Object")
}

// Refine implemnts the Shape interface
func (o *Object) Refine() []Shape {
	var shapes []Shape
	for _, mesh := range o.meshes {
		shapes = append(shapes, mesh.Refine()...)
	}
	return shapes
}
