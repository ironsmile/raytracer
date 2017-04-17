package geometry

import (
	"math"
)

const COMPARE_PRECISION = 1e-7

func Distance(p1, p2 Vector) float64 {
	return p1.Minus(p2).Length()
}

func DistanceSquared(p1, p2 Vector) float64 {
	return p1.Minus(p2).SqrLength()
}

func CoordinateSystem(vec1 Vector) (vec2 Vector, vec3 Vector) {
	if math.Abs(vec1.X) > math.Abs(vec1.Y) {
		invLen := 1.0 / math.Sqrt(vec1.X*vec1.X+vec1.Z*vec1.Z)
		vec2 = Vector{-vec1.Z * invLen, 0.0, vec1.X * invLen}
	} else {
		invLen := 1.0 / math.Sqrt(vec1.Y*vec1.Y+vec1.Z*vec1.Z)
		vec2 = Vector{0.0, vec1.Z * invLen, -vec1.Y * invLen}
	}

	vec3 = vec1.Cross(vec2)

	return
}

func Lerp(vec1, vec2 Vector, t float64) Vector {
	return vec1.MultiplyScalar(t).Plus(vec2.MultiplyScalar(1 - t))
}
