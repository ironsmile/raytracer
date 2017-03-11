package primitive

import (
	"fmt"

	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

type Rectangle struct {
	BasePrimitive
}

func (r *Rectangle) GetType() int {
	return RECTANGLE
}

func NewRectangle(w, h float64) *Rectangle {
	rec := &Rectangle{}
	rec.shape = shape.NewRectangle(w, h)
	rec.SetTransform(transform.Identity())
	return rec
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle<%s>", r.Name)
}
