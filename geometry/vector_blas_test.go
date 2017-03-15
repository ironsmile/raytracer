package geometry

import (
	"testing"

	"github.com/gonum/blas/blas64"
	"github.com/gonum/blas/cgo"
	"github.com/gonum/blas/native"
)

func TestVectorProduct(t *testing.T) {
	one := NewVector(77.345, 15.23, 2)
	other := NewVector(5, 3, 5)

	exptected := one.Product(other)

	blasOne := blas64.Vector{Inc: 1, Data: []float64{77.345, 15.23, 2}}
	blasOther := blas64.Vector{Inc: 1, Data: []float64{5, 3, 5}}

	found := blas64.Dot(len(blasOne.Data), blasOne, blasOther)

	if exptected != found {
		t.Errorf("BLAS and own imp product differs. Expected %f but found %f", exptected, found)
	}
}

func BenchmarkVectorProduct(t *testing.B) {
	one := NewVector(77.345, 15.23, 2)
	other := NewVector(5, 3, 5)

	t.Run("own impl", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			one.Product(other)
		}
	})

	blasOne := blas64.Vector{Inc: 1, Data: []float64{77.345, 15.23, 2}}
	blasOther := blas64.Vector{Inc: 1, Data: []float64{5, 3, 5}}

	blas64.Use(cgo.Implementation{})

	t.Run("cgo blas", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			blas64.Dot(len(blasOne.Data), blasOne, blasOther)
		}
	})

	blas64.Use(native.Implementation{})

	t.Run("native blas", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			blas64.Dot(len(blasOne.Data), blasOne, blasOther)
		}
	})
}
