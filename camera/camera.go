package camera

import (
	"github.com/ironsmile/raytracer/geometry"
)

var (
	DEBUG_X = 512
	DEBUG_Y = 384
)

type Camera interface {
	GenerateRay(float64, float64) (*geometry.Ray, float64)
	// GenerateRayDifferential(float64, float64) (*geometry.RayDifferential, float64)
}

type DoychoCamera struct {
	Origin    *geometry.Point
	Direction *geometry.Vector
}

func (c *DoychoCamera) GenerateRay(screenX, screenY float64) (*geometry.Ray, float64) {
	dir := geometry.Normalize((&geometry.Point{screenX, screenY, 0}).Minus(c.Origin))
	return &geometry.Ray{Origin: c.Origin, Direction: dir}, 1.0
}

func NewDoychoCamera(origin *geometry.Point) *DoychoCamera {
	return &DoychoCamera{Origin: origin, Direction: geometry.NewVector(0, 0, -1)}
}
