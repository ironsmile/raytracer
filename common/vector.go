package common

import (
	"fmt"
	"math"
)

type Vector struct {
	x, y, z float64
}

func (v *Vector) String() string {
	return fmt.Sprintf("Vector<%f, %f, %f>", v.x, v.y, v.z)
}

func (v *Vector) Multiply(other *Vector) *Vector {
	return NewVector(v.x*other.x, v.y*other.y, v.z*other.z)
}

func (v *Vector) MultiplyScalar(scalar float64) *Vector {
	return NewVector(v.x*scalar, v.y*scalar, v.z*scalar)
}

func (v *Vector) Plus(other *Vector) *Vector {
	return NewVector(v.x+other.x, v.y+other.y, v.z+other.z)
}

func (v *Vector) Minus(other *Vector) *Vector {
	return NewVector(v.x-other.x, v.y-other.y, v.z-other.z)
}

func (v *Vector) Product(other *Vector) float64 {
	return v.x*other.x + v.y*other.y + v.z*other.z
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

	return NewVector(v.y*other.z-v.z*other.y,
		v.z*other.x-v.x*other.z,
		v.x*other.y-v.y*other.x)
}

func (v *Vector) Neg() *Vector {
	return NewVector(-v.x, -v.y, -v.z)
}

func (v *Vector) Distance(other *Vector) float64 {
	x := v.x - other.x
	y := v.y - other.y
	z := v.z - other.z

	return x*x + y*y + z*z
}

func (v *Vector) Copy() *Vector {
	return NewVector(v.x, v.y, v.z)
}

func (v *Vector) Normalize() {
	l := 1.0 / v.Length()
	v.x *= l
	v.y *= l
	v.z *= l
}

func (v *Vector) Color() *Color {
	return (*Color)(v)
}

func NewVector(x, y, z float64) *Vector {
	vec := new(Vector)
	vec.x = x
	vec.y = y
	vec.z = z
	return vec
}
