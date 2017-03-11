package primitive

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

type Sphere struct {
	BasePrimitive
}

func (s *Sphere) String() string {
	sp, ok := s.shape.(*shape.Sphere)
	if !ok {
		return "Cannot type assert primitive.Sphere's shape to shape.Sphere"
	}
	return fmt.Sprintf("Sphere<center=%s, radius=%f>", sp.Center, sp.Radius)
}

func (s *Sphere) GetType() int {
	return SPHERE
}

func NewSphere(center geometry.Point, radius float64) *Sphere {
	s := &Sphere{}
	s.shape = shape.NewSphere(center, radius)
	s.Mat = *mat.NewMaterial()
	s.SetTransform(transform.Identity())
	return s
}
