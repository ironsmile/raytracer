package shape

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/utils"
)

// Sphere is centered on 0,0,0 and has a radious of 1. Use transformations to move it around
// and change its radius.
type Sphere struct {
	BasicShape

	radius float64
}

// Intersect implemnts the primitive interface
func (s *Sphere) Intersect(ray geometry.Ray, dist float64) (int, float64, geometry.Vector) {

	var d = ray.Direction
	var o = ray.Origin

	var a = d.X*d.X + d.Y*d.Y + d.Z*d.Z
	var b = 2 * (d.X*o.X + d.Y*o.Y + d.Z*o.Z)
	var c = o.X*o.X + o.Y*o.Y + o.Z*o.Z - s.radius*s.radius

	tNear, tFar, ok := utils.Quadratic(a, b, c)

	if !ok || tNear < 0 {
		return MISS, dist, ray.Direction
	}

	var retdist = tNear

	if tNear < 0 {
		retdist = tFar
	}

	if retdist > dist {
		return MISS, dist, ray.Direction
	}

	pHit := ray.Origin.Plus(ray.Direction.MultiplyScalar(retdist))

	return HIT, retdist, s.GetNormal(pHit)
}

// GetNormal implements the primitive interface
func (s *Sphere) GetNormal(pos geometry.Vector) geometry.Vector {
	return pos
}

// NewSphere returns a full sphere with a given radius
func NewSphere(rad float64) *Sphere {
	s := Sphere{radius: rad}

	s.bbox = bbox.FromPoint(geometry.NewVector(-rad, -rad, -rad))
	s.bbox = bbox.UnionPoint(s.bbox, geometry.NewVector(rad, rad, rad))

	return &s
}
