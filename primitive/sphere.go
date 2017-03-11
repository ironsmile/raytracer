package primitive

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"
	"github.com/ironsmile/raytracer/shape"
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

func (s *Sphere) Intersect(ray *geometry.Ray, dist float64) (int, float64, *geometry.Vector) {
	return s.shape.Intersect(ray, dist)
}

func NewSphere(center geometry.Point, radius float64) *Sphere {
	s := &Sphere{}
	s.shape = shape.NewSphere(center, radius)
	s.Mat = *mat.NewMaterial()
	return s
}
