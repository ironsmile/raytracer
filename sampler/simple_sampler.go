package sampler

import (
	"errors"
	"fmt"
	"image/color"
	"sync/atomic"

	"github.com/ironsmile/raytracer/film"
)

// ErrEndOfSampling would be returned by the sampler when no further sampling is needed
var ErrEndOfSampling = errors.New("End of sampling")

// SimpleSampler implements the most simple of samplers. It generates one sample per pixel
type SimpleSampler struct {
	output           film.Film
	subSamplers      []*SubSampler
	subSamplersCount uint32
	current          uint32

	stopped    bool
	continuous bool
}

// GetSubSampler ...
func (s *SimpleSampler) GetSubSampler() (ss *SubSampler, e error) {

	if s.stopped {
		e = ErrEndOfSampling
		return
	}

	sample := atomic.AddUint32(&s.current, 1) - 1

	if !s.continuous && sample >= s.subSamplersCount {
		e = ErrEndOfSampling
		return
	}

	if s.continuous && sample >= s.subSamplersCount {
		sample = sample % s.subSamplersCount
	}

	ss = s.subSamplers[sample]
	ss.Reset()

	if sample == 0 {
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
	s := &SimpleSampler{
		output: f,
	}

	const size = 32

	var countW = uint32(f.Width()) / size
	var countH = uint32(f.Height()) / size
	var count = countW * countH

	fmt.Printf("Creating %d sub samplers\n", count)

	s.subSamplersCount = count
	s.subSamplers = make([]*SubSampler, count)

	for i := uint32(0); i < count; i++ {
		sy := (i / countW) * size
		sx := (i % countW) * size

		var sw uint32 = size
		var sh uint32 = size

		if sx+sw > uint32(f.Width()) {
			sw = uint32(f.Width()) - sx
		}

		if sy+sh > uint32(f.Height()) {
			sh = uint32(f.Height()) - sy
		}

		s.subSamplers[i] = NewSubSampler(sx, sy, sw, sh, 1, s)
	}
	return s
}
