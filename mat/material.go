package mat

import (
	"github.com/ironsmile/raytracer/geometry"
)

var defaultMat = Material{
	Color: geometry.NewColor(1, 0, 0),
	Diff:  1,
}

type Material struct {
	Color *geometry.Color
	Refl  float64
	Diff  float64
	Refr  float64
}

func (m *Material) GetSpecular() float64 {
	return 1.0 - m.Diff
}

func NewMaterial() *Material {
	return &Material{Color: nil, Refl: 0.0, Diff: 0.0}
}

func DefaultMetiral() Material {
	return defaultMat
}
