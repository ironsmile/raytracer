package sampler

import (
	"errors"
	"fmt"
	"image/color"
	"sync"
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

	// pauseRequested stores 1 if there's an ongoing request for pausing the
	// sampler. If the request is cancelled (set to 0) before the next request
	// for the first sampler in the frame is made then no pausing will be initiated.
	pauseRequested atomic.Int32

	// paused stores 1 if the sampler is currently paused. It gets paused on requesting
	// the first sampler for a new frame when the `pauseRequest` has been up for a full
	// frame.
	paused atomic.Int32

	// resume is a chan on which requests for [SimpleSampler.GetSubSampler] will stay
	// stuck when the sampling is paused.
	resume chan struct{}

	// pauseFrame counts for how many frames the request for pausing has been ongoing.
	// Actually pausing will happen after a certain amount of frames. See
	// `minFamesBeforePause`.
	pauseFrame int

	// pauseLock controls the access to `resume` and `pauseFrame`.
	pauseLock *sync.RWMutex

	// minFamesBeforePause is the number of frames which a pause request must be
	// ongoing before the actual pausing. Setting a value lower than 1 will cause
	// mid-frame pauses.
	minFamesBeforePause int

	// pixList is a slice of randomly ordered pixels on the screen. Every sub sampler
	// receives a small sub-slice of this one and returns its pixels from it. This way
	// every sub sampler receives a set of pixels which are randomly distributed around
	// the screen. This way sub-samplers are known to sample every single pixel but also
	// no pixel will be repeated between samplers.
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

		if s.pauseRequested.Load() == 1 {
			s.doPause()
		}
	}

	if s.paused.Load() == 1 {
		s.pauseLock.RLock()
		ch := s.resume
		s.pauseLock.RUnlock()
		if ch != nil {
			<-ch
		}
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

// Pause requests the sampler to stop handling out sub samplers. It will do that
// at the end of the currently rendered frame. Rendering may be resumed by calling
// [SimpleSampler.Resume].
//
// If the sampler is already paused this function does nothing.
func (s *SimpleSampler) Pause() {
	if old := s.pauseRequested.Swap(1); old == 1 {
		return
	}
}

// Resume starts rendering again after a pause. If the sampler hasn't been paused
// then resume does nothing.
func (s *SimpleSampler) Resume() {
	if old := s.pauseRequested.Swap(0); old == 0 {
		return
	}
	s.paused.Store(0)

	s.pauseLock.Lock()
	defer s.pauseLock.Unlock()

	s.pauseFrame = 0
	if s.resume != nil {
		close(s.resume)
		s.resume = nil
	}
}

func (s *SimpleSampler) doPause() {
	s.pauseLock.Lock()
	defer s.pauseLock.Unlock()

	s.pauseFrame++
	if s.pauseFrame <= s.minFamesBeforePause {
		// The pause request hasn't been ongoing for long enough. Defer pausing for
		// later frame end.
		return
	}

	s.paused.Store(1)
	s.resume = make(chan struct{})
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
		output:              out,
		pauseLock:           &sync.RWMutex{},
		minFamesBeforePause: 1,
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
