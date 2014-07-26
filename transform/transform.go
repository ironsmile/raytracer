package transform

import (
	"github.com/ironsmile/raytracer/geometry"
)

type Transform struct {
	mat    *Matrix4x4
	matInv *Matrix4x4
}

func (t *Transform) Inverse() *Transform {
	return NewTransformationWihtInverse(t.matInv, t.mat)
}

func (t *Transform) IsIdentity() bool {
	return *t.mat == *NewMatrix(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1)
}

func (t *Transform) Point(point *geometry.Point) *geometry.Point {
	xp := float64(t.mat.els[0][0])*point.X + float64(t.mat.els[0][1])*point.Y +
		float64(t.mat.els[0][2])*point.Z + float64(t.mat.els[0][3])

	yp := float64(t.mat.els[1][0])*point.X + float64(t.mat.els[1][1])*point.Y +
		float64(t.mat.els[1][2])*point.Z + float64(t.mat.els[1][3])

	zp := float64(t.mat.els[2][0])*point.X + float64(t.mat.els[2][1])*point.Y +
		float64(t.mat.els[2][2])*point.Z + float64(t.mat.els[2][3])

	wp := float64(t.mat.els[3][0])*point.X + float64(t.mat.els[3][1])*point.Y +
		float64(t.mat.els[3][2])*point.Z + float64(t.mat.els[3][3])

	if wp == 1.0 {
		return geometry.NewPoint(xp, yp, zp)
	} else {
		return geometry.NewPoint(xp, yp, zp).MultiplyScalar(1.0 / wp)
	}
}

func (t *Transform) Vector(vec *geometry.Vector) *geometry.Vector {
	xp := float64(t.mat.els[0][0])*vec.X + float64(t.mat.els[0][1])*vec.Y +
		float64(t.mat.els[0][2])*vec.Z

	yp := float64(t.mat.els[1][0])*vec.X + float64(t.mat.els[1][1])*vec.Y +
		float64(t.mat.els[1][2])*vec.Z

	zp := float64(t.mat.els[2][0])*vec.X + float64(t.mat.els[2][1])*vec.Y +
		float64(t.mat.els[2][2])*vec.Z

	return geometry.NewVector(xp, yp, zp)
}

func (t *Transform) Ray(ray *geometry.Ray) *geometry.Ray {

	ret := geometry.Ray{}
	ret = *ray

	ret.Origin = t.Point(ray.Origin)
	ret.Direction = t.Vector(ray.Direction)

	return &ret
}

func (t *Transform) Multiply(other *Transform) *Transform {
	mat := t.mat.Multiply(other.mat)
	invMat := t.matInv.Multiply(other.matInv)
	return NewTransformationWihtInverse(mat, invMat)
}

func NewTransformation(mat *Matrix4x4) *Transform {
	inv, _ := mat.Inverse()
	return &Transform{mat, inv}
}

func NewTransformationWihtInverse(mat *Matrix4x4, inv *Matrix4x4) *Transform {
	return &Transform{mat, inv}
}
