package mat

import (
	"github.com/ironsmile/raytracer/geometry"
)

type Material struct {
	Color *geometry.Color
	Refl  float64
	Diff  float64
}

func (m *Material) GetSpecular() float64 {
	return 1.0 - m.Diff
}

func NewMaterial() *Material {
	return &Material{Color: nil, Refl: 0.0, Diff: 0.0}
}
