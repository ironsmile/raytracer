package scene

import (
	"fmt"
	"math"

	"github.com/ironsmile/raytracer/geometry"
)

type Sphere struct {
	BasePrimitive

	Center   *geometry.Point
	SqRadius float64
	Radius   float64
	RRadius  float64
}

func (s *Sphere) GetType() int {
	return SPHERE
}

func (s *Sphere) Intersect(ray *geometry.Ray, dist float64) (int, float64, *geometry.Vector) {
	v := ray.Origin.Minus(s.Center)
	b := -v.Product(ray.Direction)
	det := b*b - v.Product(v) + s.SqRadius

	retdist := dist
	retval := MISS
	if det <= 0 {
		return retval, retdist, nil
	}

	det = math.Sqrt(det)

	i1 := b - det
	i2 := b + det

	if i2 > 0 {
		if i1 < 0 {
			if i2 < dist {
				retdist = i2
				retval = INPRIM
			}
		} else {
			if i1 < dist {
				retdist = i1
				retval = HIT
			}
		}
	}

	intersectionPoint := ray.Origin.PlusVector(ray.Direction.MultiplyScalar(retdist))

	return retval, retdist, s.GetNormal(intersectionPoint)
}

func (s *Sphere) GetNormal(pos *geometry.Point) *geometry.Vector {
	return pos.Minus(s.Center).MultiplyScalarIP(s.RRadius)
}

func (s *Sphere) String() string {
	return fmt.Sprintf("Sphere<center=%s, radius=%f>", s.Center, s.Radius)
}

func NewSphere(center geometry.Point, radius float64) *Sphere {
	s := new(Sphere)
	s.Center = &center
	s.SqRadius = radius * radius
	s.Radius = radius
	s.RRadius = 1.0 / radius
	s.Mat = *NewMaterial()
	return s
}
