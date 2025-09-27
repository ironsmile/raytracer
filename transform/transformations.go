package transform

import (
	"math"

	"github.com/ironsmile/raytracer/geometry"
)

func Translate(delta geometry.Vector) *Transform {
	m := NewMatrix(
		1, 0, 0, (delta.X),
		0, 1, 0, (delta.Y),
		0, 0, 1, (delta.Z),
		0, 0, 0, 1)

	mInv := NewMatrix(
		1, 0, 0, -(delta.X),
		0, 1, 0, -(delta.Y),
		0, 0, 1, -(delta.Z),
		0, 0, 0, 1)

	return NewTransformationWihtInverse(m, mInv)
}

func Scale(x, y, z float64) *Transform {
	m := NewMatrix(
		(x), 0, 0, 0,
		0, (y), 0, 0,
		0, 0, (z), 0,
		0, 0, 0, 1)

	mInv := NewMatrix(
		(1.0 / x), 0, 0, 0,
		0, (1.0 / y), 0, 0,
		0, 0, (1.0 / z), 0,
		0, 0, 0, 1)

	return NewTransformationWihtInverse(m, mInv)
}

func UniformScale(s float64) *Transform {
	return Scale(s, s, s)
}

func RotateX(angle float64) *Transform {
	rad := geometry.Radians(angle)

	sin_t := math.Sin(rad)
	cos_t := math.Cos(rad)

	m := NewMatrix(
		1, 0, 0, 0,
		0, cos_t, -sin_t, 0,
		0, sin_t, cos_t, 0,
		0, 0, 0, 1)

	return NewTransformationWihtInverse(m, m.Transpose())
}

func RotateY(angle float64) *Transform {
	rad := geometry.Radians(angle)

	sin_t := math.Sin(rad)
	cos_t := math.Cos(rad)

	m := NewMatrix(
		cos_t, 0, sin_t, 0,
		0, 1, 0, 0,
		-sin_t, 0, cos_t, 0,
		0, 0, 0, 1)

	return NewTransformationWihtInverse(m, m.Transpose())
}

func RotateZ(angle float64) *Transform {
	rad := geometry.Radians(angle)

	sin_t := math.Sin(rad)
	cos_t := math.Cos(rad)

	m := NewMatrix(
		cos_t, -sin_t, 0, 0,
		sin_t, cos_t, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1)

	return NewTransformationWihtInverse(m, m.Transpose())
}

func Rotate(angle float64, a geometry.Vector) *Transform {
	a.Normalize()

	s := math.Sin(geometry.Radians(angle))
	c := math.Cos(geometry.Radians(angle))

	m := Matrix4x4{}

	m.els[0][0] = a.X*a.X + (1.0-a.X*a.X)*c
	m.els[0][1] = a.X*a.Y*(1.0-c) - a.Z*s
	m.els[0][2] = a.X*a.Z*(1.0-c) + a.Y*s
	m.els[0][3] = 0

	m.els[1][0] = a.X*a.Y*(1.0-c) + a.Z*s
	m.els[1][1] = a.Y*a.Y + (1.0-a.Y*a.Y)*c
	m.els[1][2] = a.Y*a.Z*(1.0-c) - a.X*s
	m.els[1][3] = 0

	m.els[2][0] = a.X*a.Z*(1.0-c) - a.Y*s
	m.els[2][1] = a.Y*a.Z*(1.0-c) + a.X*s
	m.els[2][2] = a.Z*a.Z + (1.0-a.Z*a.Z)*c
	m.els[2][3] = 0

	m.els[3][0] = 0
	m.els[3][1] = 0
	m.els[3][2] = 0
	m.els[3][3] = 1

	return NewTransformationWihtInverse(m, m.Transpose())
}

func LookAt(pos, look geometry.Vector, up geometry.Vector) *Transform {
	m := Matrix4x4{}

	m.els[0][3] = pos.X
	m.els[1][3] = pos.Y
	m.els[2][3] = pos.Z
	m.els[3][3] = 1

	dir := look.Minus(pos).Normalize()
	left := up.Normalize().Cross(dir).Normalize()
	newUp := dir.Cross(left)

	m.els[0][0] = left.X
	m.els[1][0] = left.Y
	m.els[2][0] = left.Z
	m.els[3][0] = 0

	m.els[0][1] = newUp.X
	m.els[1][1] = newUp.Y
	m.els[2][1] = newUp.Z
	m.els[3][1] = 0

	m.els[0][2] = dir.X
	m.els[1][2] = dir.Y
	m.els[2][2] = dir.Z
	m.els[3][2] = 0

	inv, _ := m.Inverse()

	return NewTransformationWihtInverse(inv, m)
}

func Perspective(fov, n, f float64) *Transform {
	persp := NewMatrix(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, f/(f-n), -f*n/(f-n),
		0, 0, 1, 0)

	invTanAng := 1.0 / math.Tan(geometry.Radians(fov)/2.0)

	return Scale(invTanAng, invTanAng, 1).Multiply(NewTransformation(persp))
}

func Identity() *Transform {
	ident := NewMatrix(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1)
	return NewTransformation(ident)
}
