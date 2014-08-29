package geometry

import (
	"testing"
)

func TestVectorNormalizaton(t *testing.T) {
	vec := NewVector(15, 3, 1)
	normalized := Normalize(vec)
	expected := NewVector(0.97849211, 0.19569842, 0.06523281)

	if !normalized.Equals(expected) {
		t.Errorf("Expected %s but got %s", expected, normalized)
	}

}

func TestVectorPlusFunctions(t *testing.T) {
	one := NewVector(15, 3, 1)
	other := NewVector(0.97849211, 0.19569842, 0.06523281)

	expected := one.Plus(other)
	one.PlusIP(other)

	if !one.Equals(expected) {
		t.Errorf("Expected %s but got %s", expected, one)
	}

}

/*
   Benchmarks
*/

func BenchmarkVectorPlus(t *testing.B) {
	one := NewVector(77.345, 15.23, 2)
	other := NewVector(5, 3, 5)

	for i := 0; i < t.N; i++ {
		one.Plus(other)
	}
}

func BenchmarkVectorPlusInPlace(t *testing.B) {
	one := NewVector(77.345, 15.23, 2)
	other := NewVector(5, 3, 5)

	for i := 0; i < t.N; i++ {
		one.PlusIP(other)
	}
}

func BenchmarkVectorNormalizationFunction(t *testing.B) {
	one := NewVector(77.345, 15.23, 2)

	for i := 0; i < t.N; i++ {
		Normalize(one)
	}
}

func BenchmarkVectorNormalization(t *testing.B) {
	one := NewVector(77.345, 15.23, 2)

	for i := 0; i < t.N; i++ {
		one.Normalize()
	}
}

func BenchmarkVectorNormalizationInPlace(t *testing.B) {
	one := NewVector(77.345, 15.23, 2)

	for i := 0; i < t.N; i++ {
		one.NormalizeIP()
	}
}
