package primitive

import (
	"math"
	"testing"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

func BenchmarkPrimitiveIntersection(t *testing.B) {
	ray := geometry.NewRay(
		geometry.Vector{X: 0, Y: 0, Z: -3},
		geometry.Vector{X: 0, Y: 0, Z: 1},
	)

	t.Run("Sphere", func(t *testing.B) {
		shpere := NewSphere(2)
		for i := 0; i < t.N; i++ {
			shpere.Intersect(ray)
		}
	})

	t.Run("Triangle", func(t *testing.B) {
		triangle := NewTriangle([3]geometry.Vector{
			geometry.NewVector(-1, -1, 0), // a
			geometry.NewVector(0, 1, -3),  // b
			geometry.NewVector(1, -1, 3),  // c
		})
		for i := 0; i < t.N; i++ {
			triangle.Intersect(ray)
		}
	})

	t.Run("Rectangle", func(t *testing.B) {
		rect := NewRectangle(1, 1)
		rect.SetTransform(transform.Translate(geometry.NewVector(0, 0, 30)))
		for i := 0; i < t.N; i++ {
			rect.Intersect(ray)
		}
	})
}

func TestRectangleReturnedDistanceToIntersection(t *testing.T) {
	rect := NewRectangle(1, 1)
	rect.SetTransform(transform.Translate(geometry.NewVector(0, 0, 30)))
	ray := geometry.NewRay(
		geometry.Vector{X: 0, Y: 0, Z: 0},
		geometry.Vector{X: 0, Y: 0, Z: 1},
	)

	hitMiss, distance, normal := rect.Intersect(ray)

	if hitMiss != shape.HIT {
		t.Error("The rectangle.Intersect method failed: false negative")
	}

	expectedNormal := geometry.NewVector(0, 0, -1)

	if !normal.Equals(expectedNormal) {
		t.Errorf("Wrong normal returned by rectangle.Intersect. Expected %s but got %s",
			expectedNormal, normal)
	}

	if math.Abs(distance-30) > 0.001 {
		t.Errorf("Wrong distance returned by rectangle.Intersect. Exlected %f but got %f",
			30.0, distance)
	}
}

func TestSphereIntersection(t *testing.T) {
	shpere := NewSphere(2)

	intersectRay := geometry.NewRay(
		geometry.Vector{X: 0, Y: 0, Z: -5},
		geometry.Vector{X: 0, Y: 0, Z: 1},
	)

	if hit, _, _ := shpere.Intersect(intersectRay); hit != shape.HIT {
		t.Errorf("The ray did not hit the sphere but it was expected to: %d", hit)
	}

	missRay := geometry.NewRay(
		geometry.Vector{X: 0, Y: 0, Z: -5},
		geometry.Vector{X: 0, Y: 1, Z: 0},
	)

	if hit, _, _ := shpere.Intersect(missRay); hit != shape.MISS {
		t.Errorf("The ray intersected the sphere but it was expected not to: %d", hit)
	}
}
