package sampler

import (
	"errors"
)

// ErrSubSamplerEnd represents the end this sub-sampler's cycle
var ErrSubSamplerEnd = errors.New("End of sub sampler")

// SubSampler generates samples for a rectangular subsection of a sampler
type SubSampler struct {
	x uint32
	y uint32
	w uint32

	current  uint32
	end      uint32
	perPixel uint32

	parent *SimpleSampler
}

func (s *SubSampler) GetSample() (x, y, w float64, err error) {
	if s.current >= s.end {
		err = ErrSubSamplerEnd
		return
	}
	if s.parent.stopped {
		err = ErrEndOfSampling
		return
	}
	x = float64(s.current%s.w + s.x)
	y = float64(s.current/s.w + s.y)
	w = 1
	s.current++
	return
}

func (s *SubSampler) Reset() {
	s.current = 0
}

func NewSubSampler(x, y, w, h uint32, perPixel uint32, p *SimpleSampler) *SubSampler {
	return &SubSampler{
		x:        x,
		y:        y,
		w:        w,
		perPixel: perPixel,
		end:      w * h,
		parent:   p,
	}
}
