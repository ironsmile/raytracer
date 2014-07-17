package transform

import (
	"math"
	"testing"
)

func TestIndexing(t *testing.T) {
	m := NewMatrix(
		0, 1, 2, 0,
		1, 3, 0, 1,
		2, 3, 0, 0,
		0, 0, 3, 2)

	tests := [][3]int{
		[3]int{0, 0, 0},
		[3]int{1, 1, 3},
		[3]int{2, 2, 0},
		[3]int{3, 3, 2},
		[3]int{1, 1, 3},
		[3]int{1, 2, 0},
		[3]int{0, 1, 1},
		[3]int{0, 2, 2},
		[3]int{3, 0, 0},
		[3]int{3, 2, 3},
		[3]int{3, 3, 2},
	}

	for _, test := range tests {
		expected := float32(test[2])
		i, j := test[0], test[1]

		found := m.Get(i, j)

		if found != expected {
			t.Errorf("Getting index %d,%d did not return %f but %f",
				i, j, expected, found)
		}
	}

}

func TestMultiplication(t *testing.T) {
	one := NewMatrix(
		0, 1, 2, 0,
		1, 3, 0, 1,
		2, 3, 0, 0,
		0, 0, 3, 2)

	other := NewMatrix(
		0, 0, 1, 1,
		1, 0, 1, 1,
		1, 1, 0, 0,
		0, 1, 0, 1)

	ident := NewMatrix(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1)

	result := one.Multiply(other)

	expected := NewMatrix(
		3, 2, 1, 1,
		3, 1, 4, 5,
		3, 0, 5, 5,
		3, 5, 0, 2)

	if *result != *expected {
		t.Errorf("Multiplication did not work. Expected %s but got %s", expected, result)
	}

	result = one.Multiply(ident)

	if *result != *one {
		t.Errorf("Multiplication with identity did not yield self but %s", result)
	}
}

func TestTransposition(t *testing.T) {
	one := NewMatrix(
		0, 1, 2, 0,
		1, 3, 0, 1,
		2, 3, 0, 0,
		0, 0, 3, 2)

	foundColumn := one.GetColumn(3)
	expectedColumn := [4]float32{0, 1, 0, 2}

	if foundColumn != expectedColumn {
		t.Errorf("Get column operation failed! Expected %v but got %v",
			expectedColumn, foundColumn)
	}

	expected := NewMatrix(
		0, 1, 2, 0,
		1, 3, 3, 0,
		2, 0, 0, 3,
		0, 1, 0, 2)

	transposed := one.Transpose()

	if *transposed != *expected {
		t.Errorf("Transposition failed! Expected %s but got %s", expected, transposed)
	}
}

func TestInversion(t *testing.T) {
	ident := NewMatrix(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1)

	inverted, err := ident.Inverse()

	if err != nil {
		t.Errorf("Identity has no inverse? Wrong!")
	}

	if *inverted != *ident {
		t.Errorf("Inversion of identity returned %s", inverted)
	}

	one := NewMatrix(
		0, 1, 2, 0,
		1, 3, 3, 0,
		2, 0, 0, 3,
		0, 1, 0, 2)

	inverted, err = one.Inverse()
	if err != nil {
		t.Errorf("Inverting %s returned error %s", one, err)
	}

	expected := NewMatrix(
		-0.500000, 0.333333, 0.333333, -0.500000,
		-0.666667, 0.444444, -0.222222, 0.333333,
		0.833333, -0.222222, 0.111111, -0.166667,
		0.333333, -0.222222, 0.111111, 0.333333)

	if !equal(expected, inverted) {
		t.Errorf("Expected %s but got %s", expected, inverted)
	}

}

func TestEquals(t *testing.T) {
	one := NewMatrix(
		0, 1, 2, 0,
		1, 3.345, 3, 0,
		2, 0, 0, 3,
		0, 1.23, 0, 2)

	other := NewMatrix(
		0, 1, 2, 0,
		1, 3.345, 3, 0,
		2, 0, 0, 3,
		0, 1.23, 0, 2)

	different := NewMatrix(
		0, 1, 2, 0,
		1, 3.345, 3, 0,
		2, 0, 0, 3,
		0, 1.23, 0, -1)

	epsDifferent := NewMatrix(
		0, 1, 2, 0,
		1, 3.3450001, 3, 0,
		2, 0, 0, 3,
		0, 1.23, 0, 2)

	eq := *one == *other

	if eq != true {
		t.Errorf("Equals method said %s is different from %s", one, other)
	}

	eq = *one == *different

	if eq != false {
		t.Errorf("Equals method said that %s and %s are the same", one, different)
	}

	eq = *one == *epsDifferent

	if eq != true {
		t.Errorf("Equals said that %s and %s are different", one,
			epsDifferent)
	}

}

/*
   Benchmarks
*/

func BenchmarkMultiplication(t *testing.B) {
	one := NewMatrix(
		0, 1, 2, 0,
		1, 3, 0, 1,
		2, 3, 0, 0,
		0, 0, 3, 2)

	other := NewMatrix(
		0, 0, 1, 1,
		1, 0, 1, 1,
		1, 1, 0, 0,
		0, 1, 0, 1)

	for i := 0; i < t.N; i++ {
		one.Multiply(other)
	}
}

func BenchmarkTransposition(t *testing.B) {
	one := NewMatrix(
		0, 1, 2, 0,
		1, 3, 0, 1,
		2, 3, 0, 0,
		0, 0, 3, 2)

	for i := 0; i < t.N; i++ {
		one.Transpose()
	}
}

func BenchmarkInversion(t *testing.B) {
	one := NewMatrix(
		0, 1, 2, 0,
		1, 3, 0, 1,
		2, 3, 0, 0,
		0, 0, 3, 2)

	for i := 0; i < t.N; i++ {
		one.Inverse()
	}
}

func equal(one, other *Matrix4x4) bool {
	eps := 0.000001

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if math.Abs(float64(one.els[i][j]-other.els[i][j])) > eps {
				return false
			}
		}
	}

	return true
}
