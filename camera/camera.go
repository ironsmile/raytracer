package camera

import (
	"github.com/ironsmile/raytracer/geometry"
)

var (
	// DEBUG_X = 512
	// DEBUG_Y = 384
	DEBUG_X = 0
	DEBUG_Y = 0
)

type Camera interface {
	GenerateRay(float64, float64) (*geometry.Ray, float64)
	GenerateRayIP(float64, float64, *geometry.Ray) float64
	// GenerateRayDifferential(float64, float64) (*geometry.RayDifferential, float64)

	Forward(float64) error
	Backward(float64) error
	Left(float64) error
	Right(float64) error
}
