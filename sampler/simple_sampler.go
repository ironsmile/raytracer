package sampler

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/ironsmile/raytracer/film"
)

type SimpleSampler struct {
	output film.Film

	stopChan chan bool

	samplesChan chan *Sample
	sampleFreed chan bool

	stopped bool
}

func (s *SimpleSampler) Init(f film.Film) error {
	s.output = f
	s.stopChan = make(chan bool)
	s.samplesChan = make(chan *Sample, 500)
	s.sampleFreed = make(chan bool)

	var wg sync.WaitGroup

	max := f.Width() * f.Height()
	sampleGenerators := 3
	generatorSize := max / sampleGenerators

	for i := 0; i < sampleGenerators; i++ {
		from := i * generatorSize
		to := from + generatorSize
		wg.Add(1)
		go s.sampleGenerator(from, to, &wg)
	}

	go s.cleanup(&wg)

	return nil
}

func (s *SimpleSampler) GetSample() (x float64, y float64, e error) {
	e = nil
	smpl := <-s.samplesChan
	if smpl == nil {
		e = fmt.Errorf("End of sampling")
		return
	}
	x, y = smpl.X, smpl.Y
	// s.sampleFreed <- true
	return
}

func (s *SimpleSampler) UpdateScreen(x, y float64, clr color.Color) {
	s.output.Set(int(x), int(y), clr)
}

func (s *SimpleSampler) Stop() {
	if s.stopped {
		return
	}
	close(s.stopChan)
	close(s.samplesChan)
	close(s.sampleFreed)
	s.stopped = true
}

func (s *SimpleSampler) sampleGenerator(from, to int, wg *sync.WaitGroup) {
	defer wg.Done()

	w := s.output.Width()
	var x, y float64

	for i := from; i < to; i++ {
		y = float64(i / w)
		x = float64(i % w)

		sample := &Sample{}
		sample.X, sample.Y = x, y
		s.samplesChan <- sample
	}
}

func (s *SimpleSampler) cleanup(wg *sync.WaitGroup) {
	wg.Wait()
	s.Stop()
}
