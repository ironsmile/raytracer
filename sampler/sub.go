package sampler

import (
	"errors"
	"math/rand"
)

// ErrSubSamplerEnd represents the end this sub-sampler's cycle
var ErrSubSamplerEnd = errors.New("End of sub sampler")

// SubSampler generates samples for a rectangular subsection of a sampler
type SubSampler struct {
	x uint32
	y uint32
	w uint32

	current     uint32
	end         uint32
	perPixel    uint32
	samplesDone uint32

	sampleWeight float64

	parent *SimpleSampler
}

// GetSample returns a single sample which should be raytraced
func (s *SubSampler) GetSample() (x, y, w float64, err error) {
	if s.current >= s.end {
		if s.samplesDone+1 >= s.perPixel {
			err = ErrSubSamplerEnd
			return
		}
		s.current = 0
		s.samplesDone++
	}
	if s.parent.stopped {
		err = ErrEndOfSampling
		return
	}
	x = float64(s.current%s.w+s.x) + rand.Float64()
	y = float64(s.current/s.w+s.y) + rand.Float64()
	w = s.sampleWeight
	s.current++
	return
}

// Reset returns this sub sampler to its initial condition and ready for the next frame
func (s *SubSampler) Reset() {
	s.samplesDone = 0
	s.current = 0
}

// NewSubSampler returns a sumb sampler for particular region of the screen. x and y
// represent the coordinates of this sampler in image space. w and h are the width and
// height of the sampled space respectively. perPixel dictates how many samples per
// pixle should be generated.
func NewSubSampler(x, y, w, h uint32, perPixel uint32, p *SimpleSampler) *SubSampler {
	return &SubSampler{
		x:            x,
		y:            y,
		w:            w,
		perPixel:     perPixel,
		sampleWeight: 1.0 / float64(perPixel),
		end:          w * h,
		parent:       p,
	}
}
