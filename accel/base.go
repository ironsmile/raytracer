package accel

import (
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
}

// IntersectBBoxEdge implements the Primitive interface
func (b *Base) IntersectBBoxEdge(_ geometry.Ray) bool {
	panic("IntersectBBoxEdge should not be called for accelerator")
}

// GetColor implements the primivite interface
func (b *Base) GetColor() *geometry.Color {
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
