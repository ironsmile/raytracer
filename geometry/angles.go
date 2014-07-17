package geometry

import (
	"math"
)

//  This function transforms degrees in radians
func Radians(deg float64) float64 {
	return deg * (math.Pi / 180.0)
}
