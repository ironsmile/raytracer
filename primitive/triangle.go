package primitive

import (
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

type Triangle struct {
	BasePrimitive
}

func NewTriangle(vertices [3]geometry.Vector) *Triangle {
	triangle := &Triangle{}
	triangle.shape = shape.NewTriangle(vertices)
	triangle.SetTransform(transform.Identity())
	triangle.id = GetNewID()
	return triangle
}
