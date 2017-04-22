package primitive

import (
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

func NewSphere(radius float64) *BasePrimitive {
	s := &BasePrimitive{}
	s.shape = shape.NewSphere(radius)
	s.SetTransform(transform.Identity())
	s.id = GetNewID()
	return s
}
