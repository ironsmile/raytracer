package primitive

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

type Triangle struct {
	BasePrimitive
}

func NewTriangle(vertices [3]geometry.Point) *Triangle {
	triangle := &Triangle{}
	triangle.shape = shape.NewTriangle(vertices)
	triangle.SetTransform(transform.Identity())
	return triangle
}

func (t *Triangle) String() string {
	if tr, ok := t.shape.(*shape.Triangle); ok {
		return fmt.Sprintf("Triangle<%s>: %+v", t.Name, tr.Vertices)
	}
	return "Could not type assert triangle's shape"
}

func (t *Triangle) GetType() int {
	return TRIANGLE
}
