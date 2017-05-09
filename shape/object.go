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

	embree "github.com/fogleman/go-embree"
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

	mesh      *embree.Mesh
	triangles []*MeshTriangle
}

func (o *Object) compile() {
	triangles := o.Refine()

	if o.mesh == nil {
		eTriangles := make([]embree.Triangle, len(triangles))
		o.triangles = make([]*MeshTriangle, len(triangles))

		for i, ts := range triangles {
			t := ts.(*MeshTriangle)
			v1, v2, v3 := t.getPoints()

			eTriangles[i] = embree.Triangle{
				A: embree.Vector{X: v1.X, Y: v1.Y, Z: v1.Z},
				B: embree.Vector{X: v2.X, Y: v2.Y, Z: v2.Z},
				C: embree.Vector{X: v3.X, Y: v3.Y, Z: v3.Z},
			}

			o.triangles[i] = t
		}
		o.mesh = embree.NewMesh(eTriangles)
	}
}

// Intersect implements the Shape interface
func (o *Object) Intersect(ray geometry.Ray, dg *DifferentialGeometry) bool {
	eRay := embree.Ray{
		Org: embree.Vector{X: ray.Origin.X, Y: ray.Origin.Y, Z: ray.Origin.Z},
		Dir: embree.Vector{X: ray.Direction.X, Y: ray.Direction.Y, Z: ray.Direction.Z},
	}

	hit := o.mesh.Intersect(eRay)
	if hit.Index < 0 {
		return false
	}

	if dg == nil {
		return true
	}

	dg.Distance = hit.T
	dg.Shape = o.triangles[hit.Index]

	p1, p2, p3 := o.triangles[hit.Index].getPoints()
	e1 := p2.Minus(p1)
	e2 := p3.Minus(p1)
	s1 := ray.Direction.Cross(e2)
	divisor := s1.Product(e1)
	invDivisor := 1 / divisor
	d := ray.Origin.Minus(p1)
	b1 := d.Product(s1) * invDivisor
	s2 := d.Cross(e1)
	b2 := ray.Direction.Product(s2) * invDivisor
	dg.Normal = o.triangles[hit.Index].interpolatedNormal(b1, b2)

	return true
}

// IntersectP implements the Shape interface
func (o *Object) IntersectP(ray geometry.Ray) bool {
	return o.Intersect(ray, nil)
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
					faceMath := mat.Material{
						Color: geometry.NewColor(
							foundMat.DiffuseColor.R,
							foundMat.DiffuseColor.G,
							foundMat.DiffuseColor.B,
						),
						Refr: 1 - foundMat.Dissolve,
						Diff: 1,
					}
					if faceMath.Refr > 0 {
						faceMath.RefrIndex = 1.5
					}
					faceMesh.SetMaterial(faceMath)
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

	o.compile()

	return o, nil
}

// CanIntersect implements the Shape interface
func (o *Object) CanIntersect() bool {
	return true
}

// GetMaterial implements the Shape interface
func (o *Object) GetMaterial() *mat.Material {
	panic("GetMaterial should not be called for shape.Object")
}

// SetMaterial implements Shape interface
func (o *Object) SetMaterial(mat.Material) {
	panic("SetMaterial should not be called for shape.Object")
}

// Refine implemnts the Shape interface
func (o *Object) Refine() []Shape {
	var shapes []Shape
	for _, mesh := range o.meshes {
		shapes = append(shapes, mesh.Refine()...)
	}
	return shapes
}
