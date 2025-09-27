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

	throughMiddleRay := geometry.NewRay(
		geometry.NewVector(0, 0, -2),
		geometry.NewVector(0, 0, 1),
	)

	if intersected, _, _ := box.IntersectP(throughMiddleRay); !intersected {
		t.Errorf("Ray %+v did not intersect bbox %+v", throughMiddleRay, box)
	}

	opposideDirectionRay := geometry.NewRay(
		geometry.NewVector(0, 0, -2),
		geometry.NewVector(0, 0, -1),
	)

	if intersected, _, _ := box.IntersectP(opposideDirectionRay); intersected {
		t.Errorf("Opposide direction ray %+v did intersect bbox %+v", opposideDirectionRay, box)
	}

	sideWaysRay := geometry.NewRay(
		geometry.NewVector(0, 0, -2),
		geometry.NewVector(0, 1, 0),
	)

	if intersected, _, _ := box.IntersectP(sideWaysRay); intersected {
		t.Errorf("Sideways ray %+v did intersect bbox %+v", sideWaysRay, box)
	}
}

func TestBBoxInsideMethod(t *testing.T) {
	box := New(
		geometry.NewVector(-1, -1, -1),
		geometry.NewVector(1, 1, 1),
	)

	var p = geometry.NewVector(0, 0, 0)

	if !box.Inside(p) {
		t.Errorf("bbx.Inside says that point %+v is not inside bbox %+v\n", p, box)
	}

	p = geometry.NewVector(2, 0.5, 0.5)

	if box.Inside(p) {
		t.Errorf("bbx.Inside says that point %+v is inside bbox %+v\n", p, box)
	}
}

func BenchmarkBBoxIntersections(t *testing.B) {
	box := New(
		geometry.NewVector(-1, -1, 0),
		geometry.NewVector(1, 1, 2),
	)

	throughMiddleRay := geometry.NewRay(
		geometry.NewVector(0, 0, -2),
		geometry.NewVector(0, 0, 1),
	)

	t.Run("intersected", func(t *testing.B) {
		for t.Loop() {
			box.IntersectP(throughMiddleRay)
		}
	})

	opposideDirectionRay := geometry.NewRay(
		geometry.NewVector(0, 0, -2),
		geometry.NewVector(0, 0, -1),
	)

	t.Run("oppositeDirection", func(t *testing.B) {
		for t.Loop() {
			box.IntersectP(opposideDirectionRay)
		}
	})

	sideWaysRay := geometry.NewRay(
		geometry.NewVector(0, 0, -2),
		geometry.NewVector(0, 1, 0),
	)

	t.Run("sideWays", func(t *testing.B) {
		for t.Loop() {
			box.IntersectP(sideWaysRay)
		}
	})
}
