package shape

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"
)

// The following constants are used as a return values for Intersect. HIT means that the
// ray has hit the shape, MISS - the ray has missed the shape and INPRIM means that the
// ray has hit the shape from the inside.
const (
	HIT = iota
	MISS
	INPRIM
)

// Shape is a interfece which defines a 3D shape which can be tested for intersection and stuff
type Shape interface {
	Intersect(geometry.Ray, *DifferentialGeometry) bool
	IntersectP(geometry.Ray) bool
	GetObjectBBox() *bbox.BBox
	CanIntersect() bool
	Refine() []Shape
	MaterialAt(geometry.Vector) *mat.Material
	SetMaterial(mat.Material)
	NormalAt(geometry.Vector) geometry.Vector
}

// BasicShape implements few common methods and properties among all shapes
type BasicShape struct {
	bbox     *bbox.BBox // in object space
	material *mat.Material
}

// GetObjectBBox returns a bounding box around the shape in object space or nil if no such was
// calculated.
func (b *BasicShape) GetObjectBBox() *bbox.BBox {
	return b.bbox
}

// CanIntersect implements the Shape interface
func (b *BasicShape) CanIntersect() bool {
	return true
}

// Refine implements the Shape interface
func (b *BasicShape) Refine() []Shape {
	panic("Refine should only be called on shapes which cannot be intersected: Basic")
}

// Intersect implements the Shape interface
func (b *BasicShape) Intersect(geometry.Ray, *DifferentialGeometry) bool {
	panic("Intersect is not implemented for basic shape")
}

// IntersectP implements the Shape interface
func (b *BasicShape) IntersectP(geometry.Ray) bool {
	panic("IntersectP is not implemented for basic shape")
}

// MaterialAt implements the Shape interface
func (b *BasicShape) MaterialAt(geometry.Vector) *mat.Material {
	return b.material
}

// NormalAt implements the Shape interface
func (b *BasicShape) NormalAt(geometry.Vector) geometry.Vector {
	panic("NormalAt is not implemented for basic shape")
}

// SetMaterial implements Shape interface
func (b *BasicShape) SetMaterial(mtl mat.Material) {
	if b.material == nil {
		b.material = &mtl
	} else {
		*b.material = mtl
	}
}
