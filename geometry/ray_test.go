package geometry

import (
	"testing"
)

func TestRayAtMethod(t *testing.T) {
	r := NewRay(NewVector(0, 0, 1), NewVector(0, 1, 0))

	at := r.At(2)
	expected := NewVector(0, 2, 1)

	if !at.Equals(expected) {
		t.Errorf("Expected point at distance 2 to be %+v but it was %+v\n", expected, at)
	}
}
