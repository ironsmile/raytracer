package geometry

import (
	"testing"
)

func TestLerpFunction(t *testing.T) {
	v1 := NewVector(-1, 1, 0)
	v2 := NewVector(1, 1, 0)

	expected := NewVector(0, 1, 0)
	found := Lerp(v1, v2, 0.5).Normalize()

	if !expected.Equals(found) {
		t.Errorf("Expected normal %+v but found %+v", expected, found)
	}
}
