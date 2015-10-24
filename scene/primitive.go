package scene

import (
	"github.com/ironsmile/raytracer/geometry"
)

type Primitive interface {
	GetType() int
	Intersect(*geometry.Ray, float64) (isHit int, distance float64, normal *geometry.Vector)
	GetColor() *geometry.Color
	GetMaterial() *Material
	IsLight() bool
	GetName() string
}

type BasePrimitive struct {
	Mat   Material
	Light bool
	Name  string
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

func (b *BasePrimitive) GetMaterial() *Material {
	return &b.Mat
}
