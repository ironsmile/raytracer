package sampler

import "github.com/ironsmile/raytracer/geometry"

type Sample struct {
	X, Y float64
}

type Sampler interface {
	GetSample() (float64, float64, error)
	UpdateScreen(float64, float64, *geometry.Color)
	Stop()
}
