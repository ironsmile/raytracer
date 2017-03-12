package primitive

import (
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

const (
	NOTHING = iota
	SPHERE
	PLANE
	TRIANGLE
	OBJECT
	RECTANGLE
)

type Primitive interface {
	GetType() int
	Intersect(geometry.Ray, float64) (isHit int, distance float64, normal *geometry.Vector)
	GetColor() *geometry.Color
	GetMaterial() *mat.Material
	IsLight() bool
	GetName() string
	Shape() shape.Shape
}

type BasePrimitive struct {
	Mat   mat.Material
	Light bool
	Name  string
	shape shape.Shape

	objToWorld transform.Transform
	worldToObj transform.Transform
}

func (b *BasePrimitive) GetName() string {
	return b.Name
}

func (p *BasePrimitive) IsLight() bool {
	return p.Light
}

func (b *BasePrimitive) GetColor() *geometry.Color {
	return b.Mat.Color
}

func (b *BasePrimitive) GetMaterial() *mat.Material {
	return &b.Mat
}

func (b *BasePrimitive) Intersect(ray geometry.Ray, dist float64) (int, float64, *geometry.Vector) {
	b.worldToObj.RayIP(&ray)

	res, hitDist, normal := b.shape.Intersect(&ray, dist)

	if res != shape.HIT {
		return res, hitDist, normal
	}

	return res, hitDist, b.objToWorld.NormalIP(normal)
}

func (b *BasePrimitive) Shape() shape.Shape {
	return b.shape
}

func (b *BasePrimitive) SetTransform(t *transform.Transform) {
	b.objToWorld = *t
	b.worldToObj = *t.Inverse()
}
