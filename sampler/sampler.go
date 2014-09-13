package sampler

import (
	"image/color"

	"github.com/ironsmile/raytracer/film"
)

type Sample struct {
	X, Y float64
}

type Sampler interface {
	Init(film.Film) error
	GetSample() (float64, float64, error)
	UpdateScreen(float64, float64, color.Color)
	Stop()
}
