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

// Shape is a interface which defines a 3D shape which can be tested for intersection and stuff
type Shape interface {
	// Intersect calculates the nearest intersection between a ray and this shape.
	// Returns "false" when there's no intersection at al.
	Intersect(geometry.Ray, *DifferentialGeometry) bool

	// IntersectP returns "true" if there's an intersection between the given ray
	// and _any_ geometry in the shape. Unlike [SHape.Intersect] it may not be the
	// nearest one to the ray's origin.
	IntersectP(geometry.Ray) bool

	// GetObjectBBox returns the bounding box of this shape in its object space.
	GetObjectBBox() *bbox.BBox

	// CanIntersect returns "true" if this shape can be intersected directly. Returns
	// "false" when this shape has to be refined in smaller shapes before intersection.
	CanIntersect() bool

	// Refine breaks down this shape into smaller shapes which define it. Can only be
	// called when [Shape.CanIntersect] returns "false".
	Refine() []Shape

	// MaterialAt returns the material on this shape on particular spot in world
	// space. (??? why world space?)
	MaterialAt(geometry.Vector) *mat.Material

	// SetMaterial sets a new material for this shape.
	SetMaterial(mat.Material)

	// NormalAt returns the normal on this shape on particular spot. The spot is given
	// in world space. (??? why world space?)
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
