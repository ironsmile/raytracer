package primitive

import (
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

// Cylinder is a cylinder primitive.
type Cylinder struct {
	BasePrimitive
}

// NewCylinder returns a new [Cylinder] which with a certain radius which has
// its bottom cap at `bottom` and top cap at `top`.
func NewCylinder(radius float64, bottom, top geometry.Vector) *Cylinder {
	c := &Cylinder{}
	c.shape = shape.NewCylinder(radius, bottom, top)
	c.SetTransform(transform.Identity())
	c.id = GetNewID()
	return c
}
