package geometry

import (
	"fmt"
	"math"
)

type Vector struct {
	X, Y, Z float64
}

func (v *Vector) String() string {
	return fmt.Sprintf("Vector<%.8f, %.8f, %.8f>", v.X, v.Y, v.Z)
}

func (v *Vector) Multiply(other *Vector) *Vector {
	return &Vector{v.X * other.X, v.Y * other.Y, v.Z * other.Z}
}

func (v *Vector) MultiplyIP(other *Vector) *Vector {
	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
	return v
}

func (v *Vector) MultiplyScalar(scalar float64) *Vector {
	return &Vector{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func (v *Vector) MultiplyScalarIP(scalar float64) *Vector {
	v.X *= scalar
	v.Y *= scalar
	v.Z *= scalar
	return v
}

func (v *Vector) Plus(other *Vector) *Vector {
	return &Vector{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

func (v *Vector) PlusIP(other *Vector) *Vector {
	v.X, v.Y, v.Z = v.X+other.X, v.Y+other.Y, v.Z+other.Z
	return v
}

func (v *Vector) Minus(other *Vector) *Vector {
	return &Vector{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

func (v *Vector) MinusIP(other *Vector) *Vector {
	v.X, v.Y, v.Z = v.X-other.X, v.Y-other.Y, v.Z-other.Z
	return v

}

func (v *Vector) Product(other *Vector) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v *Vector) ProductPoint(point *Point) float64 {
	return v.X*point.X + v.Y*point.Y + v.Z*point.Z
}

func (v *Vector) Length() float64 {
	return math.Sqrt(v.SqrLength())
}

func (v *Vector) SqrLength() float64 {
	return v.Product(v)
}

func (v *Vector) Dot(other *Vector) float64 {
	return v.Product(other)
}

func (v *Vector) Cross(other *Vector) *Vector {

	return &Vector{v.Y*other.Z - v.Z*other.Y,
		v.Z*other.X - v.X*other.Z,
		v.X*other.Y - v.Y*other.X}
}

func (v *Vector) CrossIP(other *Vector) *Vector {

	v.X, v.Y, v.Z = v.Y*other.Z-v.Z*other.Y,
		v.Z*other.X-v.X*other.Z,
		v.X*other.Y-v.Y*other.X
	return v
}

func (v *Vector) Neg() *Vector {
	return &Vector{-v.X, -v.Y, -v.Z}
}

func (v *Vector) NegIP() *Vector {
	v.X, v.Y, v.Z = -v.X, -v.Y, -v.Z
	return v
}

func (v *Vector) Distance(other *Vector) float64 {
	X := v.X - other.X
	Y := v.Y - other.Y
	Z := v.Z - other.Z

	return X*X + Y*Y + Z*Z
}

func (v *Vector) Copy() *Vector {
	return &Vector{v.X, v.Y, v.Z}
}

func (v *Vector) CopyToSelf(other *Vector) *Vector {
	v.X, v.Y, v.Z = other.X, other.Y, other.Z
	return v
}

func (v *Vector) Equals(other *Vector) bool {
	if math.Abs(v.X-other.X) > COMPARE_PRECISION {
		return false
	}
	if math.Abs(v.Y-other.Y) > COMPARE_PRECISION {
		return false
	}
	if math.Abs(v.Z-other.Z) > COMPARE_PRECISION {
		return false
	}
	return true
}

func (v *Vector) Normalize() *Vector {
	l := 1.0 / v.Length()
	return &Vector{v.X * l, v.Y * l, v.Z * l}
}

func (v *Vector) NormalizeIP() *Vector {
	l := 1.0 / v.Length()
	v.X *= l
	v.Y *= l
	v.Z *= l
	return v
}

func (v *Vector) Point() *Point {
	return &Point{v.X, v.Y, v.Z}
}

func (v *Vector) ByIndex(index int) float64 {
	switch index {
	case 0:
		return v.X
	case 1:
		return v.Y
	case 2:
		return v.Z
	default:
		panic("Index out of range for vector")
	}
}

func NewVector(X, Y, Z float64) *Vector {
	return &Vector{X, Y, Z}
}

func Normalize(vec *Vector) *Vector {
	return vec.Normalize()
}
