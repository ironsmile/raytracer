package bbox

import (
	"testing"

	"github.com/ironsmile/raytracer/geometry"
)

func TestBBoxIntersections(t *testing.T) {
	box := New(
		geometry.NewPoint(-1, -1, 0),
		geometry.NewPoint(1, 1, 2),
	)

	throughMiddleRay := geometry.Ray{
		Origin:    *geometry.NewPoint(0, 0, -2),
		Direction: *geometry.NewVector(0, 0, 1),
	}

	if intersected, _, _ := box.IntersectP(throughMiddleRay); !intersected {
		t.Errorf("Ray %+v did not intersect bbox %+v", throughMiddleRay, box)
	}

	opposideDirectionRay := geometry.Ray{
		Origin:    *geometry.NewPoint(0, 0, -2),
		Direction: *geometry.NewVector(0, 0, -1),
	}

	if intersected, _, _ := box.IntersectP(opposideDirectionRay); intersected {
		t.Errorf("Opposide direction ray %+v did intersect bbox %+v", opposideDirectionRay, box)
	}

	sideWaysRay := geometry.Ray{
		Origin:    *geometry.NewPoint(0, 0, -2),
		Direction: *geometry.NewVector(0, 1, 0),
	}

	if intersected, _, _ := box.IntersectP(sideWaysRay); intersected {
		t.Errorf("Sideways ray %+v did intersect bbox %+v", sideWaysRay, box)
	}
}
