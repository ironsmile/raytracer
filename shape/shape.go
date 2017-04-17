package shape

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
)

// The following constants are used as a return values for Intersect. HIT means that the
// ray has hit the shape, MISS - the ray has missed the shape and INPRIM means that the
// ray has hit the shape from the inside.
const (
	HIT = iota
	MISS
	INPRIM
)

// Shape is a interfece which defines a 3D shape which can be tested for intersection and stuff
type Shape interface {
	Intersect(geometry.Ray) (isHit int, distance float64, normal geometry.Vector)
	IntersectP(geometry.Ray) bool
	GetObjectBBox() *bbox.BBox
}

// BasicShape implements few common methods and properties among all shapes
type BasicShape struct {
	bbox *bbox.BBox // in object space
}

// GetObjectBBox returns a bounding box around the shape in object space or nil if no such was
// calculated.
func (b *BasicShape) GetObjectBBox() *bbox.BBox {
	return b.bbox
}
