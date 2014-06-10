package scene

import (
	"github.com/ironsmile/raytracer/common"
)

type Primitive interface {
	GetType() int
	Intersect(*common.Ray, float64) (int, float64)
	GetNormal(*common.Vector) *common.Vector
	GetColor() common.Color
	GetMaterial() Material
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

func (b *BasePrimitive) GetColor() common.Color {
	return *b.Mat.Color
}

func (b *BasePrimitive) GetMaterial() Material {
	return b.Mat
}
