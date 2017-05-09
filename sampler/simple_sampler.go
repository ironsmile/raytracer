package sampler

import (
	"errors"
	"image/color"
	"sync/atomic"

	"github.com/ironsmile/raytracer/film"
)

// ErrEndOfSampling would be returned by the sampler when no further sampling is needed
var ErrEndOfSampling = errors.New("End of sampling")

// SimpleSampler implements the most simple of samplers. It generates one sample per pixel
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

// GetSample returns the next (x,y) screen coordinates for whic a ray should be generated
// and traced
func (s *SimpleSampler) GetSample() (x float64, y float64, e error) {

	if s.stopped {
		e = ErrEndOfSampling
		return
	}

	sample := atomic.AddUint64(&s.currentSample, 1) - 1

	if !s.continuous && sample >= s.lastSample {
		e = ErrEndOfSampling
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

// UpdateScreen sets a pixel color for this sampler's output
func (s *SimpleSampler) UpdateScreen(x, y float64, clr color.Color) {
	s.output.Set(int(x), int(y), clr)
}

// Stop would cause all further calls to GetSample to return ErrEndOfSampling
func (s *SimpleSampler) Stop() {
	s.stopped = true
}

// MakeContinuous makes sure this sampler would contiuoue to generate samples
// in perpetuity, eventually looping back to the start of the image
func (s *SimpleSampler) MakeContinuous() {
	s.continuous = true
}

// NewSimple returns a SimpleSampler, suited for a film. This means that the sampler
// would take into consideration the film's width and height.
func NewSimple(f film.Film) *SimpleSampler {
	s := new(SimpleSampler)

	s.output = f
	s.lastSample = uint64(f.Width() * f.Height())
	s.width = uint64(f.Width())

	return s
}
