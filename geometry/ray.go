package geometry

import "math"

// Ray represents a straight line with origin and a direction
type Ray struct {
	Origin    Vector
	Direction Vector

	Mint float64
	Maxt float64

	Debug bool
}

// BackToDefaults zeroes out a ray which can then be reused somewhere int the program
func (r *Ray) BackToDefaults() {
	r.Debug = false
	r.Mint = 0
	r.Maxt = math.MaxFloat64
}

// At returns the point which as at distance t from the Origin in Direction
func (r *Ray) At(t float64) Vector {
	return r.Origin.Plus(r.Direction.MultiplyScalar(t))
}

// Intersect returns the intersection point between two rays. Two rays may not always
// intersect so that the second argument says wether there is an intersectoin at all
func (r *Ray) Intersect(o Ray) (Vector, bool) {
	n := r.Direction.Y*r.Origin.Z + r.Direction.Z*o.Origin.Y - r.Direction.Z*r.Origin.Y -
		r.Direction.Y*o.Origin.Z

	m := r.Direction.Y*o.Direction.Z - r.Direction.Z*o.Direction.Y

	v := n / m

	if v < 0 && v > -EPSILON {
		v = 0
	}

	if v < 0 || v > o.Maxt {
		return o.Direction, false
	}

	u := (o.Origin.Y + o.Direction.Y*v - r.Origin.Y) / r.Direction.Y

	if u < 0 && u > -EPSILON {
		u = 0
	}

	if u < 0 || u > r.Maxt {
		return o.Direction, false
	}

	return r.At(u), true
}

// NewRay retursn a new ray with Min zero and Max the maximum float64 value
func NewRay(origin, direction Vector) Ray {
	return Ray{
		Origin:    origin,
		Direction: direction,
		Maxt:      math.MaxFloat64,
	}
}
