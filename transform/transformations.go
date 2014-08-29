package transform

import (
	"math"

	"github.com/ironsmile/raytracer/geometry"
)

func Translate(delta *geometry.Vector) *Transform {
	m := NewMatrix(
		1, 0, 0, float32(delta.X),
		0, 1, 0, float32(delta.Y),
		0, 0, 0, float32(delta.Z),
		0, 0, 0, 1)

	mInv := NewMatrix(
		1, 0, 0, -float32(delta.X),
		0, 1, 0, -float32(delta.Y),
		0, 0, 0, -float32(delta.Z),
		0, 0, 0, 1)

	return NewTransformationWihtInverse(m, mInv)
}

func Scale(x, y, z float64) *Transform {
	m := NewMatrix(
		float32(x), 0, 0, 0,
		0, float32(y), 0, 0,
		0, 0, float32(z), 0,
		0, 0, 0, 1)

	mInv := NewMatrix(
		float32(1.0/x), 0, 0, 0,
		0, float32(1.0/y), 0, 0,
		0, 0, float32(1.0/z), 0,
		0, 0, 0, 1)

	return NewTransformationWihtInverse(m, mInv)
}

func RotateX(angle float64) *Transform {
	rad := geometry.Radians(angle)

	sin_t := float32(math.Sin(rad))
	cos_t := float32(math.Cos(rad))

	m := NewMatrix(
		1, 0, 0, 0,
		0, cos_t, -sin_t, 0,
		0, sin_t, cos_t, 0,
		0, 0, 0, 1)

	return NewTransformationWihtInverse(m, m.Transpose())
}

func RotateY(angle float64) *Transform {
	rad := geometry.Radians(angle)

	sin_t := float32(math.Sin(rad))
	cos_t := float32(math.Cos(rad))

	m := NewMatrix(
		cos_t, 0, sin_t, 0,
		0, 1, 0, 0,
		-sin_t, 0, cos_t, 0,
		0, 0, 0, 1)

	return NewTransformationWihtInverse(m, m.Transpose())
}

func RotateZ(angle float64) *Transform {
	rad := geometry.Radians(angle)

	sin_t := float32(math.Sin(rad))
	cos_t := float32(math.Cos(rad))

	m := NewMatrix(
		cos_t, -sin_t, 0, 0,
		sin_t, cos_t, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1)

	return NewTransformationWihtInverse(m, m.Transpose())
}

func Rotate(angle float64, axis *geometry.Vector) *Transform {
	a := axis.Copy()
	a.Normalize()

	s := math.Sin(geometry.Radians(angle))
	c := math.Cos(geometry.Radians(angle))

	m := &Matrix4x4{}

	m.els[0][0] = float32(a.X*a.X + (1.0-a.X*a.X)*c)
	m.els[0][1] = float32(a.X*a.Y*(1.0-c) - a.Z*s)
	m.els[0][2] = float32(a.X*a.Z*(1.0-c) + a.Y*s)
	m.els[0][3] = 0

	m.els[1][0] = float32(a.X*a.Y*(1.0-c) + a.Z*s)
	m.els[1][1] = float32(a.Y*a.Y + (1.0-a.Y*a.Y)*c)
	m.els[1][2] = float32(a.Y*a.Z*(1.0-c) - a.X*s)
	m.els[1][3] = 0

	m.els[2][0] = float32(a.X*a.Z*(1.0-c) - a.Y*s)
	m.els[2][1] = float32(a.Y*a.Z*(1.0-c) + a.X*s)
	m.els[2][2] = float32(a.Z*a.Z + (1.0-a.Z*a.Z)*c)
	m.els[2][3] = 0

	m.els[3][0] = 0
	m.els[3][1] = 0
	m.els[3][2] = 0
	m.els[3][3] = 1

	return NewTransformationWihtInverse(m, m.Transpose())
}

func LookAt(pos, look *geometry.Point, up *geometry.Vector) *Transform {
	m := &Matrix4x4{}

	m.els[0][3] = float32(pos.X)
	m.els[1][3] = float32(pos.Y)
	m.els[2][3] = float32(pos.Z)
	m.els[3][3] = 1

	dir := look.Minus(pos).NormalizeIP()
	left := up.Normalize().CrossIP(dir).NormalizeIP()
	newUp := dir.Cross(left)

	m.els[0][0] = float32(left.X)
	m.els[1][0] = float32(left.Y)
	m.els[2][0] = float32(left.Z)
	m.els[3][0] = 0

	m.els[0][1] = float32(newUp.X)
	m.els[1][1] = float32(newUp.Y)
	m.els[2][1] = float32(newUp.Z)
	m.els[3][1] = 0

	m.els[0][2] = float32(dir.X)
	m.els[1][2] = float32(dir.Y)
	m.els[2][2] = float32(dir.Z)
	m.els[3][2] = 0

	inv, _ := m.Inverse()

	return NewTransformationWihtInverse(inv, m)
}

func Perspective(fov, n, f float64) *Transform {
	persp := NewMatrix(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, float32(f/(f-n)), float32(-f*n/(f-n)),
		0, 0, 1, 0)

	invTanAng := 1.0 / math.Tan(geometry.Radians(fov)/2.0)

	return Scale(invTanAng, invTanAng, 1).Multiply(NewTransformation(persp))
}
