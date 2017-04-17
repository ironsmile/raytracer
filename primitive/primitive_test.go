package primitive

import (
	"math"
	"testing"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/transform"
)

func BenchmarkPrimitiveIntersection(t *testing.B) {
	ray := geometry.NewRay(
		geometry.Vector{X: 0, Y: 0, Z: -3},
		geometry.Vector{X: 0, Y: 0, Z: 1},
	)

	shpere := NewSphere(2)

	t.Run("Sphere.Intersect", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			shpere.Intersect(ray)
		}
	})

	t.Run("Sphere.IntersectP", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			shpere.IntersectP(ray)
		}
	})

	triangle := NewTriangle([3]geometry.Vector{
		geometry.NewVector(-1, -1, 0), // a
		geometry.NewVector(0, 1, -3),  // b
		geometry.NewVector(1, -1, 3),  // c
	})

	t.Run("Triangle.Intersect", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			triangle.Intersect(ray)
		}
	})

	t.Run("Triangle.IntersectP", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			triangle.IntersectP(ray)
		}
	})

	rect := NewRectangle(1, 1)
	rect.SetTransform(transform.Translate(geometry.NewVector(0, 0, 30)))

	t.Run("Rectangle.Intersect", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			rect.Intersect(ray)
		}
	})

	t.Run("Rectangle.IntersectP", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			rect.IntersectP(ray)
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

	pr, distance, normal := rect.Intersect(ray)

	if pr == nil {
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

func TestRectangleIntersectionWithDistance(t *testing.T) {
	rect := NewRectangle(1, 1)
	rect.SetTransform(transform.Translate(geometry.NewVector(0, 0, 30)))
	ray := geometry.NewRay(
		geometry.Vector{X: 0, Y: 0, Z: 0},
		geometry.Vector{X: 0, Y: 0, Z: 1},
	)
	ray.Mint = geometry.EPSILON
	ray.Maxt = 25

	pr, _, _ := rect.Intersect(ray)

	if pr != nil {
		t.Error("The rectangle.Intersect with maxt method failed: false positive")
	}

	if rect.IntersectP(ray) {
		t.Error("The rectangle.IntersectP with maxt method failed: false positive")
	}

	ray = geometry.NewRay(
		geometry.Vector{X: 0, Y: 0, Z: 30 - geometry.EPSILON},
		geometry.Vector{X: 0, Y: 0, Z: 1},
	)
	ray.Mint = geometry.EPSILON * 2

	pr, _, _ = rect.Intersect(ray)

	if pr != nil {
		t.Error("The rectangle.Intersect method failed: false positive")
	}

	if rect.IntersectP(ray) {
		t.Error("The rectangle.IntersectP method failed: false positive")
	}
}

func TestSphereIntersection(t *testing.T) {
	shpere := NewSphere(2)

	intersectRay := geometry.NewRay(
		geometry.Vector{X: 0, Y: 0, Z: -5},
		geometry.Vector{X: 0, Y: 0, Z: 1},
	)

	if pr, _, _ := shpere.Intersect(intersectRay); pr == nil {
		t.Errorf("The ray did not hit the sphere but it was expected to - Intersect")
	}

	if !shpere.IntersectP(intersectRay) {
		t.Errorf("The ray did not hit the sphere but it was expected to - IntersectP")
	}

	missRay := geometry.NewRay(
		geometry.Vector{X: 0, Y: 0, Z: -5},
		geometry.Vector{X: 0, Y: 1, Z: 0},
	)

	if pr, _, _ := shpere.Intersect(missRay); pr != nil {
		t.Errorf("The ray intersected the sphere but it was expected not to - Intersect")
	}

	if shpere.IntersectP(missRay) {
		t.Errorf("The ray intersected the sphere but it was expected not to - IntersectP")
	}
}
