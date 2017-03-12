package sampler

import (
	"errors"
	"image/color"
	"sync/atomic"

	"github.com/ironsmile/raytracer/film"
)

var EndOfSampling error = errors.New("End of sampling")

type SimpleSampler struct {
	output film.Film

	stopped    bool
	continuous bool

	// A counter which is used for generating samples
	currentSample uint64

	// This is the number of samples needed to generate the scene once
	lastSample uint64

	// The width of the scene. This is used when currentSample exeeds lastSample
	width uint64
}

func (s *SimpleSampler) GetSample() (x float64, y float64, e error) {

	if s.stopped {
		e = EndOfSampling
		return
	}

	sample := atomic.AddUint64(&s.currentSample, 1) - 1

	if !s.continuous && sample >= s.lastSample {
		e = EndOfSampling
		return
	}

	if s.continuous && sample >= s.lastSample {
		sample = sample % s.lastSample
	}

	y = float64(sample / s.width)
	x = float64(sample % s.width)

	if x == 0 && y == 0 {
		s.output.DoneFrame()
		s.output.StartFrame()
	}

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

func NewSimple(f film.Film) *SimpleSampler {
	s := new(SimpleSampler)

	s.output = f
	s.lastSample = uint64(f.Width() * f.Height())
	s.width = uint64(f.Width())

	return s
}
