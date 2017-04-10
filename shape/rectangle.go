package shape

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
)

type Rectangle struct {
	BasicShape

	width  float64
	height float64
}

func NewRectangle(w, h float64) *Rectangle {
	if w < 0 || w > 1 || h < 0 || h > 1 {
		panic("Recatangle width and height must be in the [0-1] region")
	}

	r := &Rectangle{width: w, height: h}

	r.bbox = bbox.FromPoint(geometry.NewVector(-0.5*w, -0.5*h, 0))
	r.bbox = bbox.UnionPoint(r.bbox, geometry.NewVector(0.5*w, 0.5*h, 0))

	return r
}

func (r *Rectangle) Intersect(ray geometry.Ray) (int, float64, geometry.Vector) {
	normal := geometry.Vector{X: 0, Y: 0, Z: -1}

	d := geometry.NewVector(0, 0, 0).Minus(ray.Origin).Product(normal)
	d /= ray.Direction.Product(normal)

	if d <= ray.Mint || d > ray.Maxt {
		return MISS, ray.Maxt, normal
	}

	hitPoint := ray.Origin.Plus(ray.Direction.MultiplyScalar(d))

	if hitPoint.X >= -0.5*r.width && hitPoint.X <= 0.5*r.width && hitPoint.Y >= -0.5*r.height && hitPoint.Y <= 0.5*r.height {
		return HIT, d, normal
	}

	return MISS, ray.Maxt, normal
}
