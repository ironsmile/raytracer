package transform

import (
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/blas/cgo"
	"github.com/ironsmile/raytracer/geometry"
)

func init() {
	blas64.Use(cgo.Implementation{})
}

func TestVectorTransform(t *testing.T) {
	vec := geometry.Vector{X: -23.55, Y: 33.77, Z: 0.032}

	t1 := NewTransformation(NewMatrix(
		0, 1.50, 2.3, 0.22,
		1, 33.2, 1.2, 1.56,
		2, 3.01, 0.1, 0.01,
		0, 0.23, 3.2, 2.12,
	))

	exptected := t1.Vector(&vec)

	blasVec := blas64.Vector{Inc: 1, Data: []float64{-23.55, 33.77, 0.032, 1}}
	blasMx := blas64.General{
		Rows:   4,
		Cols:   4,
		Stride: 4,
		Data: []float64{
			0, 1.50, 2.3, 0.22,
			1, 33.2, 1.2, 1.56,
			2, 3.01, 0.1, 0.01,
			0, 0.23, 3.2, 2.12,
		},
	}

	found := blas64.Vector{Inc: 1, Data: []float64{1, 1, 1, 1}}

	blas64.Gemv(blas.NoTrans, 1, blasMx, blasVec, 0, found)

	if exptected.X != found.Data[0] {
		t.Errorf("BLAS and own imp product differs. Expected %+v but found %+v", exptected, found)
	}
}

func BenchmarkTestVectorTransform(t *testing.B) {
	vec := geometry.Vector{X: -23.55, Y: 33.77, Z: 0.032}

	t1 := NewTransformation(NewMatrix(
		0, 1.50, 2.3, 0.22,
		1, 33.2, 1.2, 1.56,
		2, 3.01, 0.1, 0.01,
		0, 0.23, 3.2, 2.12))

	t.Run("in place", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			t1.VectorIP(&vec)
		}
	})

	blasVec := blas64.Vector{Inc: 1, Data: []float64{-23.55, 33.77, 0.032, 1}}
	blasMx := blas64.General{
		Rows:   4,
		Cols:   4,
		Stride: 4,
		Data: []float64{
			0, 1.50, 2.3, 0.22,
			1, 33.2, 1.2, 1.56,
			2, 3.01, 0.1, 0.01,
			0, 0.23, 3.2, 2.12,
		},
	}

	found := blas64.Vector{Inc: 1, Data: []float64{1, 1, 1, 1}}

	t.Run("blas", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			blas64.Gemv(blas.NoTrans, 1, blasMx, blasVec, 0, found)
		}
	})
}
