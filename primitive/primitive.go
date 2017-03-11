package primitive

import (
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"
	"github.com/ironsmile/raytracer/shape"
)

const (
	NOTHING = iota
	SPHERE
	PLANE
	TRIANGLE
	OBJECT
)

type Primitive interface {
	GetType() int
	Intersect(*geometry.Ray, float64) (isHit int, distance float64, normal *geometry.Vector)
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

func (b *BasePrimitive) Intersect(r *geometry.Ray, d float64) (int, float64, *geometry.Vector) {
	return b.shape.Intersect(r, d)
}

func (b *BasePrimitive) Shape() shape.Shape {
	return b.shape
}
