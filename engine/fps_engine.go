package engine

import (
	"fmt"
	"sync"
	"time"

	"github.com/ironsmile/raytracer/camera"
	"github.com/ironsmile/raytracer/geometry"
)

type FPSEngine struct {
	Engine
	wg       sync.WaitGroup
	stopChan chan bool
}

func (e *FPSEngine) Render() {

	quads := 3
	quadWidth := e.Width / quads
	quadHeight := e.Height / quads

	e.stopChan = make(chan bool)

	e.Dest.StartFrame()

	for quadIndX := 0; quadIndX < quads; quadIndX++ {
		for quadIndY := 0; quadIndY < quads; quadIndY++ {

			quadXStart := quadIndX * quadWidth
			quadXStop := quadXStart + quadWidth - 1

			quadYStart := quadIndY * quadHeight
			quadYStop := quadYStart + quadHeight - 1

			e.wg.Add(1)
			go e.subRender(quadXStart, quadXStop, quadYStart, quadYStop, &e.wg)
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

func (e *FPSEngine) subRender(startX, stopX, startY, stopY int,
	wg *sync.WaitGroup) {
	defer wg.Done()

	ray := &geometry.Ray{}
	accColor := geometry.NewColor(0, 0, 0)
	for {
		select {
		case _ = <-e.stopChan:
			return
		default:
			for y := startY; y <= stopY; y++ {
				for x := startX; x <= stopX; x++ {

					weight := e.Camera.GenerateRayIP(float64(x), float64(y), ray)

					if x == camera.DEBUG_X && y == camera.DEBUG_Y {
						fmt.Printf("Debugging ray:\n%v\n", ray)
						ray.Debug = true
					}

					e.Raytrace(ray, 1, accColor)

					e.Dest.Set(x, y, accColor.MultiplyScalarIP(weight))

				}
			}
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
