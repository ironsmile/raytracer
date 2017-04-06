package utils

import (
	"math"
)

// Lerp does a linear interpolation between two points
func Lerp(t, v1, v2 float64) float64 {
	return (1.0-t)*v1 + t*v2
}

// Clamp clamps the given value val to be between the values low and high
func Clamp(val, low, high float64) float64 {
	if val < low {
		return low
	}

	if val > high {
		return high
	}

	return val
}

func ConcentricSampleDisk(u1, u2 float64) (float64, float64) {
	var r, theta float64

	sx := 2*u1 - 1
	sy := 2*u2 - 1

	if sx == 0.0 && sy == 0.0 {
		return 0.0, 0.0
	}

	if sx >= -sy {
		if sx > sy {
			r = sx
			if sy > 0.0 {
				theta = sy / r
			} else {
				theta = 8.0 + sy/r
			}
		} else {
			r = sy
			theta = 2.0 - sx/r
		}
	} else {
		if sx <= sy {
			r = -sx
			theta = 4.0 - sy/r
		} else {
			r = -sy
			theta = 6.0 + sx/r
		}
	}
	theta *= math.Pi / 4.0

	return r * math.Cos(theta), r * math.Sin(theta)
}

// Min returns the minimal value betwee two float64s
func Min(a, b float64) float64 {
	if a <= b {
		return a
	}
	return b
}

// Max returns the maximal value betwee two float64s
func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Quadratic solves a quadratic equation and returns the two solutions of there are any.
// Its last return value is a boolean and true when there is a solution. The first two
// values are the solutions.
func Quadratic(a, b, c float64) (float64, float64, bool) {
	discrim := b*b - 4*a*c
	if discrim <= 0 {
		return 0, 0, false
	}
	rootDiscrim := math.Sqrt(discrim)
	var q float64
	if b < 0 {
		q = -0.5 * (b - rootDiscrim)
	} else {
		q = -0.5 * (b + rootDiscrim)
	}

	t0, t1 := q/a, c/q

	if t0 > t1 {
		t0, t1 = t1, t0
	}

	return t0, t1, true
}
