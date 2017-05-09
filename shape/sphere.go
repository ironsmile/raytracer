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

// Intersect implemnts the shape interface
func (s *Sphere) Intersect(ray geometry.Ray, dg *DifferentialGeometry) bool {

	var d = ray.Direction
	var o = ray.Origin

	var a = d.X*d.X + d.Y*d.Y + d.Z*d.Z
	var b = 2 * (d.X*o.X + d.Y*o.Y + d.Z*o.Z)
	var c = o.X*o.X + o.Y*o.Y + o.Z*o.Z - s.radius*s.radius

	tNear, tFar, ok := utils.Quadratic(a, b, c)

	if !ok || tNear < 0 {
		return false
	}

	var retdist = tNear

	if tNear < 0 {
		retdist = tFar
	}

	if retdist > ray.Maxt || retdist < ray.Mint {
		return false
	}

	if dg == nil {
		return true
	}

	dg.Shape = s
	dg.Distance = retdist
	dg.Normal = s.GetNormal(ray.At(retdist))

	return true
}

// IntersectP implements the shape interface
func (s *Sphere) IntersectP(ray geometry.Ray) bool {
	return s.Intersect(ray, nil)
}

// GetNormal implements the primitive interface
func (s *Sphere) GetNormal(pos geometry.Vector) geometry.Vector {
	return pos.Normalize()
}

// NewSphere returns a full sphere with a given radius
func NewSphere(rad float64) *Sphere {
	s := Sphere{radius: rad}

	s.bbox = bbox.FromPoint(geometry.NewVector(-rad, -rad, -rad))
	s.bbox = bbox.UnionPoint(s.bbox, geometry.NewVector(rad, rad, rad))

	return &s
}
