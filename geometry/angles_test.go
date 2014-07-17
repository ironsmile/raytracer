package geometry

import (
	"math"
	"testing"
)

func TestRadiansFunction(t *testing.T) {
	rad := Radians(180)

	if rad != math.Pi {
		t.Errorf("Expected %f but got %f for 180 degrees", math.Pi, rad)
	}

	rad = Radians(90)
	expected := math.Pi / 2

	if rad != expected {
		t.Errorf("Expected %f but got %f for 90 degrees", expected, rad)
	}
}

func BenchmarkRadiansFunction(t *testing.B) {
	for i := 0; i < t.N; i++ {
		Radians(73)
	}
}
