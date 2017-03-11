package shape

import (
	"github.com/ironsmile/raytracer/geometry"
)

const (
	HIT = iota
	MISS
	INPRIM
)

// Shape is a interfece which defines a 3D shape which can be tested for intersection and stuff
type Shape interface {
	Intersect(*geometry.Ray, float64) (isHit int, distance float64, normal *geometry.Vector)
	// GetNormal(*geometry.Point) *geometry.Vector
}
