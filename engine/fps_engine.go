package engine

// import (
// 	"sync"
// 	"time"
// )

// type FPSEngine struct {
// 	Engine
// 	wg       sync.WaitGroup
// 	stopChan chan bool
// }

// func (e *FPSEngine) Render() {

// 	e.stopChan = make(chan bool)

// 	e.Dest.StartFrame()

// 	e.startParallelRendering(&e.wg, e.renderContinuously)

// 	e.wg.Add(1)
// 	go e.screenRefresher()
// }

// func (e *FPSEngine) screenRefresher() {
// 	defer e.wg.Done()

// 	for {
// 		select {
// 		case _ = <-e.stopChan:
// 			return
// 		default:
// 			e.Dest.StartFrame()
// 			time.Sleep(20 * time.Millisecond)
// 			e.Dest.DoneFrame()
// 		}
// 	}
// }

// func (e *FPSEngine) renderContinuously(startX, stopX, startY, stopY int,
// 	wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	for {
// 		select {
// 		case _ = <-e.stopChan:
// 			return
// 		default:
// 			e.subRender(startX, stopX, startY, stopY)
// 		}

// 	}
// }

// func (e *FPSEngine) StopRendering() {
// 	close(e.stopChan)
// 	e.wg.Wait()
// 	e.Dest.Wait()
// }

// func NewFPSEngine() *FPSEngine {
// 	eng := new(FPSEngine)
// 	return eng
// }
