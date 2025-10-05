package sampler

import (
	"errors"
	"math/rand"
	"time"
)

// ErrSubSamplerEnd represents the end this sub-sampler's cycle
var ErrSubSamplerEnd = errors.New("End of sub sampler")

// SubSampler generates samples for a rectangular subsection of a sampler
type SubSampler struct {
	pixArray []sampledPixel

	current     uint32
	perPixel    uint32
	samplesDone uint32

	parent *SimpleSampler

	rand *rand.Rand
}

// GetSample returns a single sample which should be raytraced
func (s *SubSampler) GetSample() (x, y float64, err error) {
	if s.current >= uint32(len(s.pixArray)) {
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
	sample := s.pixArray[s.current]

	x = float64(sample.x) + s.rand.Float64()
	y = float64(sample.y) + s.rand.Float64()
	s.current++
	return
}

// Reset returns this sub sampler to its initial condition and ready for the next frame
func (s *SubSampler) Reset() {
	s.samplesDone = 0
	s.current = 0
}

// NewSubSampler returns a sumb sampler.... TODO
func NewSubSampler(pixArray []sampledPixel, perPixel uint32, p *SimpleSampler) *SubSampler {
	src := rand.NewSource(time.Now().UnixMicro())
	rnd := rand.New(src)

	return &SubSampler{
		pixArray: pixArray,
		perPixel: perPixel,
		parent:   p,
		rand:     rnd,
	}
}
