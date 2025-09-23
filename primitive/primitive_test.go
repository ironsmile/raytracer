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

	in := Intersection{}
	shpere := NewSphere(2)

	t.Run("Sphere.Intersect", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			shpere.Intersect(ray, &in)
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
			triangle.Intersect(ray, &in)
		}
	})

	t.Run("Triangle.IntersectP", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			triangle.IntersectP(ray)
		}
	})

	rect := NewQuad(
		geometry.NewVector(-1, 1, 0),
		geometry.NewVector(1, 1, 0),
		geometry.NewVector(1, -1, 0),
		geometry.NewVector(-1, -1, 0),
	)
	rect.SetTransform(transform.Translate(geometry.NewVector(0, 0, 30)))

	t.Run("Rectangle.Intersect", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			rect.Intersect(ray, &in)
		}
	})

	t.Run("Rectangle.IntersectP", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			rect.IntersectP(ray)
		}
	})
}

func TestRectangleReturnedDistanceToIntersection(t *testing.T) {
	rect := NewQuad(
		geometry.NewVector(-1, 1, 0),
		geometry.NewVector(1, 1, 0),
		geometry.NewVector(1, -1, 0),
		geometry.NewVector(-1, -1, 0),
	)
	rect.SetTransform(transform.Translate(geometry.NewVector(0, 0, 30)))
	ray := geometry.NewRay(
		geometry.Vector{X: 0, Y: 0, Z: 0},
		geometry.Vector{X: 0, Y: 0, Z: 1},
	)

	in := Intersection{}

	if !rect.Intersect(ray, &in) {
		t.Error("The rectangle.Intersect method failed: false negative")
	}

	expectedNormal := geometry.NewVector(0, 0, -1)
	pi := ray.At(in.DfGeometry.Distance)
	actualNormal := in.DfGeometry.Shape.NormalAt(pi)

	if !actualNormal.Equals(expectedNormal) {
		t.Errorf("Wrong normal returned by rectangle.Intersect. Expected %s but got %s",
			expectedNormal, actualNormal)
	}

	if math.Abs(in.DfGeometry.Distance-30) > 0.001 {
		t.Errorf("Wrong distance returned by rectangle.Intersect. Exlected %f but got %f",
			30.0, in.DfGeometry.Distance)
	}
}

func TestRectangleIntersectionWithDistance(t *testing.T) {
	rect := NewQuad(
		geometry.NewVector(-1, 1, 0),
		geometry.NewVector(1, 1, 0),
		geometry.NewVector(1, -1, 0),
		geometry.NewVector(-1, -1, 0),
	)
	rect.SetTransform(transform.Translate(geometry.NewVector(0, 0, 30)))
	ray := geometry.NewRay(
		geometry.Vector{X: 0, Y: 0, Z: 0},
		geometry.Vector{X: 0, Y: 0, Z: 1},
	)
	ray.Mint = geometry.EPSILON
	ray.Maxt = 25

	in := Intersection{}

	if rect.Intersect(ray, &in) {
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

	if rect.Intersect(ray, &in) {
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

	in := Intersection{}

	if !shpere.Intersect(intersectRay, &in) {
		t.Errorf("The ray did not hit the sphere but it was expected to - Intersect")
	}

	if !shpere.IntersectP(intersectRay) {
		t.Errorf("The ray did not hit the sphere but it was expected to - IntersectP")
	}

	missRay := geometry.NewRay(
		geometry.Vector{X: 0, Y: 0, Z: -5},
		geometry.Vector{X: 0, Y: 1, Z: 0},
	)

	if shpere.Intersect(missRay, &in) {
		t.Errorf("The ray intersected the sphere but it was expected not to - Intersect")
	}

	if shpere.IntersectP(missRay) {
		t.Errorf("The ray intersected the sphere but it was expected not to - IntersectP")
	}
}
