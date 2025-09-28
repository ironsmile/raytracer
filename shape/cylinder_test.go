package shape_test

import (
	"fmt"
	"testing"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/shape"
)

// TestCylinderNormals checks the normals calculations for a cylinder.
func TestCylinderNormals(t *testing.T) {
	cyl := shape.NewCylinder(
		1, geometry.NewVector(0, 0, 0), geometry.NewVector(0, 2, 0),
	)

	tests := []struct {
		at       geometry.Vector
		expected geometry.Vector
	}{
		{
			at:       geometry.NewVector(1, 1, 0),
			expected: geometry.NewVector(1, 0, 0),
		},
		{
			at:       geometry.NewVector(-1, 1, 0),
			expected: geometry.NewVector(-1, 0, 0),
		},
		{
			at:       geometry.NewVector(0, 1, 1),
			expected: geometry.NewVector(0, 0, 1),
		},
		{
			at:       geometry.NewVector(0, 1, -1),
			expected: geometry.NewVector(0, 0, -1),
		},
		{
			at:       geometry.NewVector(0, 0, -1),
			expected: geometry.NewVector(0, 0, -1),
		},
		{
			at:       geometry.NewVector(0, 0, 1),
			expected: geometry.NewVector(0, 0, 1),
		},
		{
			at:       geometry.NewVector(1, 0, 0),
			expected: geometry.NewVector(1, 0, 0),
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("check_%d", i), func(t *testing.T) {
			actual := cyl.NormalAt(test.at)

			if !test.expected.Equals(actual) {
				t.Fatalf("expected normal %s but got %s\n", test.expected, actual)
			}
		})
	}
}
