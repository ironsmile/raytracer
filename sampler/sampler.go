package sampler

import (
	"image/color"
)

type Sample struct {
	X, Y float64
}

type Sampler interface {
	GetSample() (float64, float64, error)
	UpdateScreen(float64, float64, color.Color)
	Stop()
}
