package primitive

import (
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

func NewRectangle(w, h float64) *BasePrimitive {
	rec := &BasePrimitive{}
	rec.shape = shape.NewRectangle(w, h)
	rec.SetTransform(transform.Identity())
	return rec
}
