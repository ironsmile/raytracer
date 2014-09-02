package engine

import (
	"sync"
	"time"
)

type FPSEngine struct {
	Engine
	wg       sync.WaitGroup
	stopChan chan bool
}

func (e *FPSEngine) Render() {

	e.stopChan = make(chan bool)

	e.Dest.StartFrame()

	quads := 3
	quadWidth := e.Width / quads
	quadHeight := e.Height / quads

	for quadIndX := 0; quadIndX < quads; quadIndX++ {
		for quadIndY := 0; quadIndY < quads; quadIndY++ {

			quadXStart := quadIndX * quadWidth
			quadXStop := quadXStart + quadWidth - 1

			quadYStart := quadIndY * quadHeight
			quadYStop := quadYStart + quadHeight - 1

			e.wg.Add(1)
			go e.startSubRender(quadXStart, quadXStop, quadYStart, quadYStop, &e.wg)
		}
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
			e.Dest.StartFrame()
			time.Sleep(100 * time.Millisecond)
			e.Dest.DoneFrame()
		}
	}
}

func (e *FPSEngine) startSubRender(startX, stopX, startY, stopY int,
	wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case _ = <-e.stopChan:
			return
		default:
			e.subRender(startX, stopX, startY, stopY)
		}

	}
}

func (e *FPSEngine) StopRendering() {
	close(e.stopChan)
	e.wg.Wait()
	e.Dest.Wait()
}

func NewFPSEngine() *FPSEngine {
	eng := new(FPSEngine)
	return eng
}
