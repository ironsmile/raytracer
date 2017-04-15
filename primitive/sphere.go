package primitive

import (
	"fmt"

	"github.com/ironsmile/raytracer/mat"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

type Sphere struct {
	BasePrimitive
}

func (s *Sphere) String() string {
	return fmt.Sprintf("Sphere<transofrm=%s>", s.worldToObj)
}

func NewSphere(radius float64) *Sphere {
	s := &Sphere{}
	s.shape = shape.NewSphere(radius)
	s.Mat = *mat.NewMaterial()
	s.SetTransform(transform.Identity())
	return s
}
