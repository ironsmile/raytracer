package transform

import (
	"testing"
)

func TestMatrixIndexing(t *testing.T) {
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
		expected := float64(test[2])
		i, j := test[0], test[1]

		found := m.Get(i, j)

		if found != expected {
			t.Errorf("Getting index %d,%d did not return %f but %f",
				i, j, expected, found)
		}
	}

}

func TestMatrixMultiplication(t *testing.T) {
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

	if !result.Equals(expected) {
		t.Errorf("Multiplication did not work. Expected %s but got %s", expected, result)
	}

	result = one.Multiply(ident)

	if !result.Equals(one) {
		t.Errorf("Multiplication with identity did not yield self but %s", result)
	}
}

func TestMatrixTransposition(t *testing.T) {
	one := NewMatrix(
		0, 1, 2, 0,
		1, 3, 0, 1,
		2, 3, 0, 0,
		0, 0, 3, 2)

	foundColumn := one.GetColumn(3)
	expectedColumn := [4]float64{0, 1, 0, 2}

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

	if !transposed.Equals(expected) {
		t.Errorf("Transposition failed! Expected %s but got %s", expected, transposed)
	}
}

func TestMatrixInversion(t *testing.T) {
	ident := NewMatrix(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1)

	inverted, err := ident.Inverse()

	if err != nil {
		t.Errorf("Identity has no inverse? Wrong!")
	}

	if !inverted.Equals(ident) {
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
		-0.5000000, 0.3333333, 0.3333333, -0.5000000,
		-0.6666667, 0.4444444, -0.2222222, 0.3333333,
		0.8333333, -0.2222222, 0.1111111, -0.1666667,
		0.3333333, -0.2222222, 0.1111111, 0.3333333)

	if !expected.Equals(inverted) {
		t.Errorf("Expected %s but got %s", expected, inverted)
	}

}

func TestMatrixEquals(t *testing.T) {
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

	eq := one.Equals(other)

	if eq != true {
		t.Errorf("Equals method said %s is different from %s", one, other)
	}

	eq = one.Equals(different)

	if eq != false {
		t.Errorf("Equals method said that %s and %s are the same", one, different)
	}

	eq = one.Equals(epsDifferent)

	if eq != true {
		t.Errorf("Equals said that %s and %s are different", one,
			epsDifferent)
	}

}

/*
   Benchmarks
*/

func BenchmarkMatrixMultiplication(t *testing.B) {
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

func BenchmarkMatrixTransposition(t *testing.B) {
	one := NewMatrix(
		0, 1, 2, 0,
		1, 3, 0, 1,
		2, 3, 0, 0,
		0, 0, 3, 2)

	for i := 0; i < t.N; i++ {
		one.Transpose()
	}
}

func BenchmarkMatrixInversion(t *testing.B) {
	one := NewMatrix(
		0, 1, 2, 0,
		1, 3, 0, 1,
		2, 3, 0, 0,
		0, 0, 3, 2)

	for i := 0; i < t.N; i++ {
		one.Inverse()
	}
}
