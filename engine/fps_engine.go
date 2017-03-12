package engine

import (
	"runtime"
	"sync"
	"time"

	"github.com/ironsmile/raytracer/sampler"
)

type FPSEngine struct {
	Engine
	wg       sync.WaitGroup
	stopChan chan bool
}

func (e *FPSEngine) Render() {

	e.stopChan = make(chan bool)

	e.Dest.StartFrame()

	for i := 0; i < runtime.NumCPU(); i++ {
		e.wg.Add(1)
		go e.subRender(&e.wg)
	}

	e.wg.Add(1)
	go e.screenRefresher()
}

func (e *FPSEngine) screenRefresher() {
	defer e.wg.Done()

	for {
		select {
		case _ = <-e.stopChan:
			return
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (e *FPSEngine) StopRendering() {
	close(e.stopChan)
	e.wg.Wait()
	e.Dest.Wait()
}

func NewFPS(smpl sampler.Sampler) *FPSEngine {
	eng := new(FPSEngine)
	initEngine(&eng.Engine, smpl)
	return eng
}
