package shape

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
)

// Rectangle represents a rectangular shape in the 3d space. It is implemented as a
// rectangle, centered at the 0,0 point in the xy plane and then rotated with
// transformations.
type Rectangle struct {
	BasicShape

	width  float64
	height float64
}

// NewRectangle returns a pointer to a rectangle with width `w` and height `h`
func NewRectangle(w, h float64) *Rectangle {
	if w < 0 || w > 1 || h < 0 || h > 1 {
		panic("Recatangle width and height must be in the [0-1] region")
	}

	r := &Rectangle{width: w, height: h}

	r.bbox = bbox.FromPoint(geometry.NewVector(-0.5*w, -0.5*h, 0))
	r.bbox = bbox.UnionPoint(r.bbox, geometry.NewVector(0.5*w, 0.5*h, 0))

	return r
}

// Intersect implements the Shape interface
func (r *Rectangle) Intersect(ray geometry.Ray, dg *DifferentialGeometry) bool {
	normal := geometry.Vector{X: 0, Y: 0, Z: -1}

	d := geometry.NewVector(0, 0, 0).Minus(ray.Origin).Product(normal)
	d /= ray.Direction.Product(normal)

	if d < ray.Mint || d > ray.Maxt {
		return false
	}

	hp := ray.At(d)

	if hp.X < -0.5*r.width || hp.X > 0.5*r.width || hp.Y < -0.5*r.height || hp.Y > 0.5*r.height {
		return false
	}

	dg.Distance = d
	dg.Normal = normal

	return true
}

// IntersectP implements the Shape interface
func (r *Rectangle) IntersectP(ray geometry.Ray) bool {
	normal := geometry.Vector{X: 0, Y: 0, Z: -1}

	d := geometry.NewVector(0, 0, 0).Minus(ray.Origin).Product(normal)
	d /= ray.Direction.Product(normal)

	if d < ray.Mint || d > ray.Maxt {
		return false
	}

	hp := ray.At(d)

	if hp.X >= -0.5*r.width && hp.X <= 0.5*r.width && hp.Y >= -0.5*r.height && hp.Y <= 0.5*r.height {
		return true
	}

	return false
}
