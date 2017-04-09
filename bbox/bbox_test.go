package bbox

import (
	"testing"

	"github.com/ironsmile/raytracer/geometry"
)

func TestBBoxIntersections(t *testing.T) {
	box := New(
		geometry.NewVector(-1, -1, 0),
		geometry.NewVector(1, 1, 2),
	)

	throughMiddleRay := geometry.Ray{
		Origin:    geometry.NewVector(0, 0, -2),
		Direction: geometry.NewVector(0, 0, 1),
	}

	if intersected, _, _ := box.IntersectP(throughMiddleRay); !intersected {
		t.Errorf("Ray %+v did not intersect bbox %+v", throughMiddleRay, box)
	}

	opposideDirectionRay := geometry.Ray{
		Origin:    geometry.NewVector(0, 0, -2),
		Direction: geometry.NewVector(0, 0, -1),
	}

	if intersected, _, _ := box.IntersectP(opposideDirectionRay); intersected {
		t.Errorf("Opposide direction ray %+v did intersect bbox %+v", opposideDirectionRay, box)
	}

	sideWaysRay := geometry.Ray{
		Origin:    geometry.NewVector(0, 0, -2),
		Direction: geometry.NewVector(0, 1, 0),
	}

	if intersected, _, _ := box.IntersectP(sideWaysRay); intersected {
		t.Errorf("Sideways ray %+v did intersect bbox %+v", sideWaysRay, box)
	}
}

func BenchmarkBBoxIntersections(t *testing.B) {
	box := New(
		geometry.NewVector(-1, -1, 0),
		geometry.NewVector(1, 1, 2),
	)

	throughMiddleRay := geometry.Ray{
		Origin:    geometry.NewVector(0, 0, -2),
		Direction: geometry.NewVector(0, 0, 1),
	}

	t.Run("intersected", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			box.IntersectP(throughMiddleRay)
		}
	})

	opposideDirectionRay := geometry.Ray{
		Origin:    geometry.NewVector(0, 0, -2),
		Direction: geometry.NewVector(0, 0, -1),
	}

	t.Run("oppositeDirection", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			box.IntersectP(opposideDirectionRay)
		}
	})

	sideWaysRay := geometry.Ray{
		Origin:    geometry.NewVector(0, 0, -2),
		Direction: geometry.NewVector(0, 1, 0),
	}

	t.Run("sideWays", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			box.IntersectP(sideWaysRay)
		}
	})
}
