package utils

import (
	"math"
)

func Lerp(t, v1, v2 float64) float64 {
	return (1.0-t)*v1 + t*v2
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

func Min(a, b float64) float64 {
	if a <= b {
		return a
	}
	return b
}

func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
