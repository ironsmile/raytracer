package camera

import (
	"github.com/ironsmile/raytracer/geometry"
)

type Camera struct {
	width  int
	height int

	Origin    *geometry.Point
	Direction *geometry.Vector
	Distance  float64

	// projectionPlane *geometry.Vector
}

func (c *Camera) Set(origin *geometry.Point) {
	c.Origin = origin
	c.Direction = geometry.NewVector(0, 0, -1)

	// c.projectionPlane = c.Origin.Plus(c.Direction.MultiplyScalar(c.Distance))
}

func (c *Camera) GetWorldPosition(screenX, screenY float64) *geometry.Vector {
	return (&geometry.Point{screenX, screenY, 0}).Minus(c.Origin)
}

func NewCamera(width, height int) *Camera {
	cam := new(Camera)
	cam.width = width
	cam.height = height
	cam.Distance = 5
	return cam
}
