package geometry

import (
	"math"
)

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

// NewRay retursn a new ray with Min zero and Max the maximum float64 value
func NewRay(origin, direction Vector) Ray {
	return Ray{
		Origin:    origin,
		Direction: direction,
		Maxt:      math.MaxFloat64,
	}
}
