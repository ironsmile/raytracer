package primitive

import (
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

type Quad struct {
	BasePrimitive
}

func NewQuad(v1, v2, v3, v4 geometry.Vector) *Quad {
	q := &Quad{}
	q.shape = shape.NewQuad(v1, v2, v3, v4)
	q.SetTransform(transform.Identity())
	q.id = GetNewID()
	return q
}
