package geometry

import (
	"testing"
)

func BenchmarkChainOfCommonOperations(t *testing.B) {
	normal := Vector{X: 0, Y: 0, Z: -1}
	ray := Ray{
		Origin:    Vector{X: 5, Y: 4, Z: -1},
		Direction: Vector{X: 0, Y: 2, Z: -1},
	}

	for i := 0; i < t.N; i++ {
		d := NewVector(0, 0, 0).Minus(ray.Origin).Product(normal)
		d /= ray.Direction.Product(normal)
	}
}
