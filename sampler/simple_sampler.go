package sampler

import (
	"fmt"
	"image/color"
	"sync/atomic"

	"github.com/ironsmile/raytracer/film"
)

type SimpleSampler struct {
	output film.Film

	stopped    bool
	continuous bool

	currentSample uint64
	lastSample    uint64
	width         uint64
}

func (s *SimpleSampler) Init(f film.Film) error {
	s.output = f

	s.lastSample = uint64(f.Width() * f.Height())
	s.width = uint64(f.Width())

	return nil
}

func (s *SimpleSampler) GetSample() (x float64, y float64, e error) {

	if s.stopped {
		e = fmt.Errorf("End of sampling")
		return
	}

	sample := atomic.AddUint64(&s.currentSample, 1) - 1

	if !s.continuous && sample >= s.lastSample {
		e = fmt.Errorf("End of sampling")
		return
	}

	if s.continuous && sample >= s.lastSample {
		sample = sample % s.lastSample
	}

	y = float64(sample / s.width)
	x = float64(sample % s.width)

	return
}

func (s *SimpleSampler) UpdateScreen(x, y float64, clr color.Color) {
	s.output.Set(int(x), int(y), clr)
}

func (s *SimpleSampler) Stop() {
	if s.stopped {
		return
	}
	s.stopped = true
}

func (s *SimpleSampler) MakeContinuous() {
	s.continuous = true
}
