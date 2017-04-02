package shape

import (
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/utils"
)

// Sphere is centered on 0,0,0 and has a radious of 1. Use transformations to move it around
// and change its radius.
type Sphere struct {
	BasicShape

	Radius float64
}

// Intersect implemnts the primitive interface
func (s *Sphere) Intersect(ray geometry.Ray, dist float64) (int, float64, geometry.Vector) {

	var d = ray.Direction
	var o = ray.Origin

	var a = d.X*d.X + d.Y*d.Y + d.Z*d.Z
	var b = 2 * (d.X*o.X + d.Y*o.Y + d.Z*o.Z)
	var c = o.X*o.X + o.Y*o.Y + o.Z*o.Z - s.Radius*s.Radius

	tNear, tFar, ok := utils.Quadratic(a, b, c)

	if !ok || tNear < 0 {
		return MISS, dist, ray.Direction
	}

	var retdist = tNear

	if tNear < 0 {
		retdist = tFar
	}

	intersectionPoint := ray.Origin.PlusVector(ray.Direction.MultiplyScalar(retdist))

	return HIT, retdist, *s.GetNormal(intersectionPoint)
}

// GetNormal implements the primitive interface
func (s *Sphere) GetNormal(pos *geometry.Point) *geometry.Vector {
	return pos.Vector()
}

// NewSphere returns a sphere
func NewSphere(radius float64) *Sphere {
	return &Sphere{Radius: radius}
}
