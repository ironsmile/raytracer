package primitive

import (
	"math"
	"testing"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

func BenchmarkPrimitiveIntersection(t *testing.B) {
	ray := geometry.Ray{
		Origin:    geometry.Point{X: 0, Y: 0, Z: -3},
		Direction: geometry.Vector{X: 0, Y: 0, Z: 1},
	}

	t.Run("Sphere", func(t *testing.B) {
		shpere := NewSphere(geometry.Point{X: 0, Y: 0, Z: 0}, 2)
		for i := 0; i < t.N; i++ {
			shpere.Intersect(ray, 1000000.0)
		}
	})

	t.Run("Triangle", func(t *testing.B) {
		triangle := NewTriangle([3]geometry.Point{
			*geometry.NewPoint(-1, -1, 0), // a
			*geometry.NewPoint(0, 1, -3),  // b
			*geometry.NewPoint(1, -1, 3),  // c
		})
		for i := 0; i < t.N; i++ {
			triangle.Intersect(ray, 1000000.0)
		}
	})

	t.Run("Rectangle", func(t *testing.B) {
		rect := NewRectangle(1, 1)
		rect.SetTransform(transform.Translate(geometry.NewVector(0, 0, 30)))
		for i := 0; i < t.N; i++ {
			rect.Intersect(ray, 1000000.0)
		}
	})
}

func TestRectangleReturnedDistanceToIntersection(t *testing.T) {
	rect := NewRectangle(1, 1)
	rect.SetTransform(transform.Translate(geometry.NewVector(0, 0, 30)))
	ray := geometry.Ray{
		Origin:    geometry.Point{X: 0, Y: 0, Z: 0},
		Direction: geometry.Vector{X: 0, Y: 0, Z: 1},
	}

	hitMiss, distance, normal := rect.Intersect(ray, 100)

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
