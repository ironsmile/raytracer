package camera

import (
	"github.com/ironsmile/raytracer/geometry"
)

type Camera interface {
	GenerateRay(float64, float64) (*geometry.Ray, float64)
	// GenerateRayDifferential(float64, float64) (*geometry.RayDifferential, float64)
}

type BasicCamera struct {
	Origin    *geometry.Point
	Direction *geometry.Vector
}

func (c *BasicCamera) SetOrigin(origin *geometry.Point) {
	c.Origin = origin
	c.Direction = geometry.NewVector(0, 0, -1)
}

func (c *BasicCamera) GenerateRay(screenX, screenY float64) (*geometry.Ray, float64) {
	dir := (&geometry.Point{screenX, screenY, 0}).Minus(c.Origin)
	dir.Normalize()
	return &geometry.Ray{Origin: c.Origin, Direction: dir}, 1.0
}

func NewBasicCamera() *BasicCamera {
	return &BasicCamera{}
}
