package primitive

import (
	"github.com/ironsmile/raytracer/mat"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

func NewSphere(radius float64) *BasePrimitive {
	s := &BasePrimitive{}
	s.shape = shape.NewSphere(radius)
	s.Mat = *mat.NewMaterial()
	s.SetTransform(transform.Identity())
	return s
}
