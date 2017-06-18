package accel

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/color"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"
	"github.com/ironsmile/raytracer/primitive"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

// Base is a accelerator which implements all the stuff that a accelerator needs to
// implement the Primitive interface. And some shared data
type Base struct {
	primitives []primitive.Primitive
	bounds     *bbox.BBox
}

// GetWorldBBox implements the Primitive interface
func (b *Base) GetWorldBBox() *bbox.BBox {
	return b.bounds
}

// IntersectBBoxEdge implements the Primitive interface
func (b *Base) IntersectBBoxEdge(ray geometry.Ray) bool {
	in, _ := b.GetWorldBBox().IntersectEdge(ray)
	return in
}

// CanIntersect implements the Primitive interface. And sice the only purpose of a grid accelerator
// is to be itnersectable, it returns true in all cases.
func (b *Base) CanIntersect() bool {
	return true
}

// Refine implements the Primitive interface
func (b *Base) Refine() []primitive.Primitive {
	panic("Refine should not be called for accelerator")
}

// GetColor implements the primivite interface
func (b *Base) GetColor() *color.Color {
	panic("GetColor should not be called for accelerator")
}

// GetMaterial implements the primivite interface
func (b *Base) GetMaterial() *mat.Material {
	panic("GetMaterial should not be called for accelerator")
}

// IsLight implements the primivite interface
func (b *Base) IsLight() bool {
	panic("IsLight should not be called for accelerator")
}

// GetLightSource implements the primivite interface
func (b *Base) GetLightSource() geometry.Vector {
	panic("GetLightSource should not be called for accelerator")
}

// GetName implements the primivite interface
func (b *Base) GetName() string {
	return "GridAccel"
}

// Shape implements the primivite interface
func (b *Base) Shape() shape.Shape {
	panic("Shape should not be called for accelerator")
}

// GetTransforms implements the primitive interface
func (b *Base) GetTransforms() (*transform.Transform, *transform.Transform) {
	panic("GetTransforms should not be called for accelerator")
}

// SetTransform implements the primitive interface
func (b *Base) SetTransform(*transform.Transform) {
	panic("SetTransform should not be called for accelerator")
}

// GetID implements the primitive interface
func (b *Base) GetID() uint64 {
	panic("GetID should not be called for accelerator")
}
