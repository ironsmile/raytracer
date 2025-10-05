package sampler

import (
	"errors"
	"fmt"
	"image/color"
	"sync/atomic"
)

// ErrEndOfSampling would be returned by the sampler when no further sampling is needed
var ErrEndOfSampling = errors.New("End of sampling")

// SimpleSampler implements the most simple of samplers. It generates a fixed amount of
// sample per pixel
type SimpleSampler struct {
	output           Output
	subSamplers      []*SubSampler
	subSamplersCount uint32
	current          uint32

	stopped    bool
	continuous bool

	pixList []sampledPixel
}

// GetSubSampler returns a rectangular sampler for a smaller section of the screen.
func (s *SimpleSampler) GetSubSampler() (*SubSampler, error) {

	if s.stopped {
		return nil, ErrEndOfSampling
	}

	sample := atomic.AddUint32(&s.current, 1) - 1

	if !s.continuous && sample >= s.subSamplersCount {
		return nil, ErrEndOfSampling
	}

	if s.continuous && sample >= s.subSamplersCount {
		sample = sample % s.subSamplersCount
	}

	ss := s.subSamplers[sample]
	ss.Reset()

	if sample == 0 {
		s.output.DoneFrame()
		s.output.StartFrame()
	}

	return ss, nil
}

// UpdateScreen sets a pixel color for this sampler's output
func (s *SimpleSampler) UpdateScreen(x, y float64, clr color.Color) {
	s.output.Set(int(x), int(y), clr)
}

// Stop would cause all further calls to GetSample to return ErrEndOfSampling
func (s *SimpleSampler) Stop() {
	s.stopped = true
}

// MakeContinuous makes sure this sampler would continue to generate samples
// in perpetuity, eventually looping back to the start of the image
func (s *SimpleSampler) MakeContinuous() {
	s.continuous = true
}

// NewSimple returns a SimpleSampler which would generate samples for a 2D output
// with certain width and height.
func NewSimple(width, height int, out Output) *SimpleSampler {
	s := &SimpleSampler{
		output: out,
	}

	sampleSet := make(map[sampledPixel]struct{})
	for x := range uint32(width) {
		for y := range uint32(height) {
			sampleSet[sampledPixel{x: x, y: y}] = struct{}{}
		}
	}

	s.pixList = make([]sampledPixel, 0, width*height)

	// Take advantage of the fact that iterating over a map returns its keys in
	// a random order.
	for k := range sampleSet {
		s.pixList = append(s.pixList, k)
	}

	var splits = uint32(width / 32)
	var count = splits * splits
	var perSampler = uint32((width * height) / int(count))

	fmt.Printf("Creating %d sub samplers\n", count)

	s.subSamplersCount = count
	s.subSamplers = make([]*SubSampler, count)

	for i := range count {
		start := i * perSampler
		end := start + perSampler
		if end > uint32(len(s.pixList)) {
			end = uint32(len(s.pixList))
		}

		subPixels := s.pixList[start:end]
		s.subSamplers[i] = NewSubSampler(subPixels, 4, s)
	}
	return s
}

type sampledPixel struct {
	x, y uint32
}
