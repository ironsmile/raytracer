package shape

import (
	"fmt"
	"math"

	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
)

// Cylinder represents a finite cylinder with a particular radius.
type Cylinder struct {
	BasicShape

	radius float64

	endcapTop    geometry.Vector
	endcapBottom geometry.Vector
}

// NewCylinder returns a new [Cylinder] which with a certain radius which has
// its bottom cap at `bottom` and top cap at `top`.
func NewCylinder(radius float64, bottom, top geometry.Vector) *Cylinder {
	c := &Cylinder{
		radius:       radius,
		endcapTop:    top,
		endcapBottom: bottom,
	}

	c.bbox = bbox.FromPoint(bottom)
	c.bbox = bbox.UnionPoint(c.bbox, top)

	spPoints := [6]geometry.Vector{
		bottom.Plus(geometry.NewVector(1, 0, 0).MultiplyScalar(radius)),
		bottom.Plus(geometry.NewVector(-1, 0, 0).MultiplyScalar(radius)),
		bottom.Plus(geometry.NewVector(0, 1, 0).MultiplyScalar(radius)),
		bottom.Plus(geometry.NewVector(0, -1, 0).MultiplyScalar(radius)),
		bottom.Plus(geometry.NewVector(0, 0, 1).MultiplyScalar(radius)),
		bottom.Plus(geometry.NewVector(0, 0, -1).MultiplyScalar(radius)),
	}

	for _, p := range spPoints {
		c.bbox = bbox.UnionPoint(c.bbox, p)
	}

	return c
}

// Intersect implements the Shape interface for the cylinder using the
// David J. Cobb method described at
// https://davidjcobb.github.io/articles/ray-cylinder-intersection.
func (c *Cylinder) Intersect(ray geometry.Ray, dg *DifferentialGeometry) bool {
	var tDist float64

	Rl := ray.Origin.Minus(c.endcapBottom)
	Cs := c.endcapTop.Minus(c.endcapBottom)
	Ch := Cs.Length()
	Ca := Cs.MultiplyScalar(1 / Ch)

	CaDotRd := Ca.Dot(ray.Direction)
	CaDotRl := Ca.Dot(Rl)
	RlDotRl := Rl.Dot(Rl)

	ca := 1 - (CaDotRd * CaDotRd)
	cb := 2 * (ray.Direction.Dot(Rl) - CaDotRd*CaDotRl)
	cc := RlDotRl - CaDotRl*CaDotRl - (c.radius * c.radius)

	hitNear, hitAway, count := quadraticRoots(ca, cb, cc)
	if count == 0 {
		// There is no intersection between a line (i.e. a "double-sided" ray) and the
		// infinite cylinder that matches our finite cylinder. This means that we cannot
		// be hitting any part of the cylinder: if we were hitting the base from the
		// inside, for example, then the "back of our ray" would be hitting the upper
		// part of the infinite cylinder.
		return false
	}

	if count > 2 {
		panic("quadratic equation with more than two answers?")
	}

	var (
		valid1 = true
		valid2 = true
	)

	Hp1 := ray.Origin.Plus(ray.Direction.MultiplyScalar(hitNear))
	Hp2 := ray.Origin.Plus(ray.Direction.MultiplyScalar(hitAway))
	Ho1 := c.endcapTop.Minus(Hp1).Dot(Ca)
	Ho2 := c.endcapTop.Minus(Hp2).Dot(Ca)

	validCount := count
	if hitNear < 0 || Ho1 < 1.0e-7 || Ho1 > Ch {
		valid1 = false
		validCount--
	}
	if hitAway < 0 || Ho2 < 1.0e-7 || Ho2 > Ch {
		valid2 = false
		if count > 1 {
			validCount--
		}
	}
	if validCount == 0 {
		// The ray never hits the bounded cylinder's curved surface. If we're looking
		// along the cylinder's axis -- whether from inside or outside -- then the ray
		// could still hit an endcap.
		//
		// Ignore cap hits at the moment. This is not what the original code of
		// David J. Cobb does. But so far it hasn't been needed so I've skipped it.
		return false
	}

	if dg == nil {
		return true
	}

	if validCount == 1 {
		if valid1 {
			tDist = hitNear
		} else if valid2 {
			tDist = hitAway
		} else {
			return false
		}
	} else {
		tDist = math.Min(hitNear, hitAway)
	}

	dg.Shape = c
	dg.Distance = tDist

	return true
}

// NormalAt implements the Shape interface
func (c *Cylinder) NormalAt(at geometry.Vector) geometry.Vector {
	dir := c.endcapTop.Minus(c.endcapBottom)
	vat := at.Minus(c.endcapBottom)
	th := math.Acos(dir.Dot(vat) / (dir.Length() * vat.Length()))
	if math.IsNaN(th) {
		fmt.Printf("NaN acos in normal calculations for cylinder at %#v\n", at)
		return at.Minus(c.endcapBottom)
	}

	vatProjLen := vat.Length() * math.Cos(th)
	axisAt := dir.Normalize().MultiplyScalar(vatProjLen)

	return at.Minus(axisAt).Normalize()
}

// IntersectP implements the Shape interface
func (c *Cylinder) IntersectP(ray geometry.Ray) bool {
	return c.Intersect(ray, nil)
}

// Given coefficients in a quadratic equation, this function gives you the roots
// and returns the number of roots. If there is only one root, then both root
// variables are set to the same value.
func quadraticRoots(a, b, c float64) (float64, float64, int) {
	const EPSILON = 1e-7

	discr := (b * b) - (4.0 * a * c)
	if discr > EPSILON {
		var bTerm float64
		if b < EPSILON {
			bTerm = -b + math.Sqrt(discr)
		} else {
			bTerm = -b - math.Sqrt(discr)
		}

		lower := bTerm / (2.0 * a) // quadratic formula
		upper := (2.0 * c) / bTerm // citardauq formula

		if lower > upper {
			lower, upper = upper, lower
		}

		return lower, upper, 2
	} else if discr > -EPSILON && discr <= EPSILON {
		lower := -(b / 2.0 * a)
		return lower, lower, 1
	}

	return 0, 0, 0
}
