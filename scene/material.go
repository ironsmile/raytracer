package scene

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
	mat := new(Material)
	col := geometry.NewColor(0.2, 0.2, 0.2)
	mat.Color = col
	mat.Refl = 0.0
	mat.Diff = 0.2

	return mat
}
