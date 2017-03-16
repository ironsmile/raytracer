package primitive

import (
	"testing"

	"github.com/ironsmile/raytracer/geometry"
)

func BenchmarkPrimitiveIntersection(t *testing.B) {
	prim := NewSphere(geometry.Point{X: 0, Y: 0, Z: 0}, 2)
	ray := geometry.Ray{
		Origin:    geometry.Point{X: 0, Y: 0, Z: -3},
		Direction: geometry.Vector{X: 0, Y: 0, Z: 1},
	}

	for i := 0; i < t.N; i++ {
		prim.Intersect(ray, 1000000.0)
	}
}
