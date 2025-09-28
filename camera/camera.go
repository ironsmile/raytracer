package camera

import (
	"github.com/ironsmile/raytracer/geometry"
)

// Camera is the interface all types of cameras have to implement
type Camera interface {
	GenerateRay(float64, float64) geometry.Ray
	// GenerateRayIP(float64, float64, *geometry.Ray) float64
	// GenerateRayDifferential(float64, float64) (*geometry.RayDifferential, float64)

	Forward(float64) error
	Backward(float64) error
	Left(float64) error
	Right(float64) error
	Up(float64) error
	Down(float64) error

	Yaw(float64) error
	Pitch(float64) error
}
