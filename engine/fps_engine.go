package engine

import (
	"runtime"
	"sync"

	"github.com/ironsmile/raytracer/sampler"
)

// FPSEngine is a type of [Engine] which is suitable for running in a window
// as part of an GUI application.
type FPSEngine struct {
	Engine
	wg sync.WaitGroup
}

// Render starts the engine and rendering.
func (e *FPSEngine) Render() {
	e.Dest.StartFrame()

	for i := 0; i < runtime.NumCPU(); i++ {
		e.wg.Add(1)
		go e.subRender(&e.wg)
	}
}

// StopRendering waits for the rendering to stop.
func (e *FPSEngine) StopRendering() {
	e.wg.Wait()
	e.Dest.Wait()
}

// Pause pauses rendering. Rendering may be continued with [FPSEngine.Resume].
func (e *FPSEngine) Pause() {
	e.Sampler.Pause()
}

// Resume stats rendering after a previous [FPSEngine.Pause].
func (e *FPSEngine) Resume() {
	e.Sampler.Resume()
}

// NewFPS returns a new FPS engine which will use the given sampler for getting
// samples for rendering.
func NewFPS(smpl *sampler.SimpleSampler) *FPSEngine {
	eng := new(FPSEngine)
	initEngine(&eng.Engine, smpl)
	return eng
}
