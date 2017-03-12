package transform

import (
	"fmt"

	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
)

type Transform struct {
	mat    Matrix4x4
	matInv Matrix4x4
}

func (t *Transform) Inverse() *Transform {
	return NewTransformationWihtInverse(&t.matInv, &t.mat)
}

func (t *Transform) IsIdentity() bool {
	return t.mat.Equals(NewMatrix(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1))
}

func (t *Transform) Point(point *geometry.Point) *geometry.Point {
	p := geometry.NewPoint(
		float64(t.mat.els[0][0])*point.X+float64(t.mat.els[0][1])*point.Y+
			float64(t.mat.els[0][2])*point.Z+float64(t.mat.els[0][3]),

		float64(t.mat.els[1][0])*point.X+float64(t.mat.els[1][1])*point.Y+
			float64(t.mat.els[1][2])*point.Z+float64(t.mat.els[1][3]),

		float64(t.mat.els[2][0])*point.X+float64(t.mat.els[2][1])*point.Y+
			float64(t.mat.els[2][2])*point.Z+float64(t.mat.els[2][3]),
	)

	wp := float64(t.mat.els[3][0])*point.X + float64(t.mat.els[3][1])*point.Y +
		float64(t.mat.els[3][2])*point.Z + float64(t.mat.els[3][3])

	if wp != 1.0 {
		p.MultiplyScalarIP(1.0 / wp)
	}
	return p
}

func (t *Transform) PointIP(point *geometry.Point) *geometry.Point {
	xp := float64(t.mat.els[0][0])*point.X + float64(t.mat.els[0][1])*point.Y +
		float64(t.mat.els[0][2])*point.Z + float64(t.mat.els[0][3])

	yp := float64(t.mat.els[1][0])*point.X + float64(t.mat.els[1][1])*point.Y +
		float64(t.mat.els[1][2])*point.Z + float64(t.mat.els[1][3])

	zp := float64(t.mat.els[2][0])*point.X + float64(t.mat.els[2][1])*point.Y +
		float64(t.mat.els[2][2])*point.Z + float64(t.mat.els[2][3])

	wp := float64(t.mat.els[3][0])*point.X + float64(t.mat.els[3][1])*point.Y +
		float64(t.mat.els[3][2])*point.Z + float64(t.mat.els[3][3])

	point.X, point.Y, point.Z = xp, yp, zp

	if wp != 1.0 {
		point.MultiplyScalarIP(1.0 / wp)
	}
	return point
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

func (t *Transform) VectorIP(vec *geometry.Vector) *geometry.Vector {
	xp := float64(t.mat.els[0][0])*vec.X + float64(t.mat.els[0][1])*vec.Y +
		float64(t.mat.els[0][2])*vec.Z

	yp := float64(t.mat.els[1][0])*vec.X + float64(t.mat.els[1][1])*vec.Y +
		float64(t.mat.els[1][2])*vec.Z

	zp := float64(t.mat.els[2][0])*vec.X + float64(t.mat.els[2][1])*vec.Y +
		float64(t.mat.els[2][2])*vec.Z

	vec.X, vec.Y, vec.Z = xp, yp, zp
	return vec
}

func (t *Transform) Ray(ray *geometry.Ray) *geometry.Ray {
	ret := *ray

	ret.Origin = *t.Point(&ray.Origin)
	ret.Direction = *t.Vector(&ray.Direction)

	return &ret
}

func (t *Transform) RayIP(ray *geometry.Ray) *geometry.Ray {

	t.PointIP(&ray.Origin)
	t.VectorIP(&ray.Direction)

	return ray
}

func (t *Transform) Normal(vec *geometry.Vector) *geometry.Vector {
	return t.NormalIP(vec.Copy())
}

func (t *Transform) NormalIP(vec *geometry.Vector) *geometry.Vector {
	var x, y, z = vec.X, vec.Y, vec.Z

	vec.X = float64(t.matInv.els[0][0])*x + float64(t.matInv.els[1][0])*y + float64(t.matInv.els[2][0])*z
	vec.Y = float64(t.matInv.els[0][1])*x + float64(t.matInv.els[1][1])*y + float64(t.matInv.els[2][1])*z
	vec.Z = float64(t.matInv.els[0][2])*x + float64(t.matInv.els[1][2])*y + float64(t.matInv.els[2][2])*z

	return vec.NormalizeIP()
}

func (t *Transform) BBox(bb *bbox.BBox) *bbox.BBox {
	ret := bbox.FromPoint(t.Point(geometry.NewPoint(bb.Min.X, bb.Min.Y, bb.Min.Z)))
	ret = bbox.UnionPoint(ret, t.Point(geometry.NewPoint(bb.Max.X, bb.Min.Y, bb.Min.Z)))
	ret = bbox.UnionPoint(ret, t.Point(geometry.NewPoint(bb.Min.X, bb.Max.Y, bb.Min.Z)))
	ret = bbox.UnionPoint(ret, t.Point(geometry.NewPoint(bb.Min.X, bb.Min.Y, bb.Max.Z)))
	ret = bbox.UnionPoint(ret, t.Point(geometry.NewPoint(bb.Min.X, bb.Max.Y, bb.Max.Z)))
	ret = bbox.UnionPoint(ret, t.Point(geometry.NewPoint(bb.Max.X, bb.Max.Y, bb.Min.Z)))
	ret = bbox.UnionPoint(ret, t.Point(geometry.NewPoint(bb.Max.X, bb.Min.Y, bb.Max.Z)))
	ret = bbox.UnionPoint(ret, t.Point(geometry.NewPoint(bb.Max.X, bb.Max.Y, bb.Max.Z)))
	return ret
}

func (t *Transform) Multiply(other *Transform) *Transform {
	mat := t.mat.Multiply(&other.mat)
	invMat := other.matInv.Multiply(&t.matInv)
	return NewTransformationWihtInverse(mat, invMat)
}

func (t *Transform) Equals(other *Transform) bool {
	return t.mat.Equals(&other.mat) && t.matInv.Equals(&other.matInv)
}

func (t *Transform) String() string {
	return fmt.Sprintf("Transformation with %+v", t.mat)
}

func (t *Transform) SwapsHandedness() bool {
	m := t.mat.els
	det := ((m[0][0] *
		(m[1][1]*m[2][2] -
			m[1][2]*m[2][1])) -
		(m[0][1] *
			(m[1][0]*m[2][2] -
				m[1][2]*m[2][0])) +
		(m[0][2] *
			(m[1][0]*m[2][1] -
				m[1][1]*m[2][0])))
	return det < 0.0
}

func NewTransformation(mat *Matrix4x4) *Transform {
	inv, _ := mat.Inverse()
	return &Transform{*mat, *inv}
}

func NewTransformationWihtInverse(mat *Matrix4x4, inv *Matrix4x4) *Transform {
	return &Transform{*mat, *inv}
}
