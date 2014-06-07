package engine

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ironsmile/raytracer/common"
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/scene"
)

const (
	EPSION = 0.00001
)

type Engine struct {
	WX1, WY1, WX2, WY2, DiffX, DiffY, SX, SY float64
	Scene                                    *scene.Scene
	Dest                                     film.Film
	Width, Height                            int64
}

func (e *Engine) SetTarget(target film.Film) {
	e.Width = int64(target.Width())
	e.Height = int64(target.Height())
	e.Dest = target
}

func (e *Engine) InitRender() {
	e.WX1 = -4
	e.WX2 = 4
	e.WY1 = 3
	e.SY = 3
	e.WY2 = -3

	e.DiffX = (e.WX2 - e.WX1) / float64(e.Width)
	e.DiffY = (e.WY2 - e.WY1) / float64(e.Height)

	fmt.Printf("Engine initialized with viewport %dx%d\n", e.Width, e.Height)
}

func (e *Engine) Raytrace(ray *common.Ray, depth int64) (
	scene.Primitive, float64, *common.Color) {

	retColor := common.NewColor(0, 0, 0)

	if depth > common.TRACEDEPTH {
		return nil, 0, retColor
	}

	prim, retdist := e.Scene.Intersect(ray)

	if prim == nil {
		return nil, 0, retColor
	}

	if prim.IsLight() {
		clr := prim.GetColor()
		retColor = &clr
		return prim, retdist, retColor
	}

	primMat := prim.GetMaterial()

	pi := ray.Origin.Plus(ray.Direction.MultiplyScalar(retdist))

	for l := 0; l < e.Scene.GetNrLights(); l++ {
		light := e.Scene.GetLight(l)
		shade := 1.0

		L := (light.(*scene.Sphere)).Center.Minus(pi)
		L.Normalize()

		if light.GetType() == scene.SPHERE {
			piOffset := pi.Plus(L.MultiplyScalar(EPSION))
			shadowRay := common.NewRay(*piOffset, *L)
			// shadowRay.Debug = ray.Debug

			intersected, _ := e.Scene.Intersect(shadowRay)

			if light != intersected {
				shade = 0.0
			}
		}

		N := prim.GetNormal(pi)

		if primMat.Diff > 0 {
			dot := N.Product(L)
			if dot > 0 {
				weight := dot * primMat.Diff * shade
				retColor = retColor.Vector().Plus(light.GetMaterial().Color.Vector().
					Multiply(primMat.Color.Vector()).MultiplyScalar(weight)).Color()
			}
		}

		if primMat.GetSpecular() > 0 {
			V := ray.Direction
			R := L.Minus(N.MultiplyScalar(2.0 * L.Product(N)))
			dot := V.Product(R)
			if dot > 0 {
				spec := math.Pow(dot, 20) * primMat.GetSpecular() * shade
				retColor = retColor.Vector().Plus(light.GetMaterial().Color.
					Vector().MultiplyScalar(spec)).Color()
			}
		}

	}

	// Reflection
	if primMat.Refl > 0.0 {
		N := prim.GetNormal(pi)
		R := ray.Direction.Minus(N.MultiplyScalar(ray.Direction.Product(N) * 2.0))

		refRay := common.NewRay(*pi.Plus(R.MultiplyScalar(EPSION)), *R)
		// refRay.Debug = ray.Debug
		_, _, refColor := e.Raytrace(refRay, depth+1)

		retColor = retColor.Vector().Plus(primMat.Color.Vector().Multiply(
			refColor.Vector()).MultiplyScalar(primMat.Refl)).Color()
	}

	return prim, retdist, retColor
}

func (e *Engine) Render() bool {

	refreshTime := 30 * time.Millisecond
	timer := time.NewTimer(refreshTime)
	timerStop := make(chan bool)

	go func() {
		defer timer.Stop()
		for {
			select {
			case _ = <-timer.C:
				e.Dest.Ping()
				timer.Reset(refreshTime)
			case _ = <-timerStop:
				return
			}
		}
	}()

	origin := common.NewVector(0, 0, -5)

	quads := 16
	quadWidth := int(e.Width) / quads
	quadHeight := int(e.Height) / quads

	var wg sync.WaitGroup

	for quadIndX := 0; quadIndX < quads; quadIndX++ {
		for quadIndY := 0; quadIndY < quads; quadIndY++ {

			quadXStart := quadIndX * quadWidth
			quadXStop := quadXStart + quadWidth - 1

			quadYStart := quadIndY * quadHeight
			quadYStop := quadYStart + quadHeight - 1

			wg.Add(1)
			go e.subRender(origin, quadXStart, quadXStop, quadYStart, quadYStop, &wg)
		}
	}

	wg.Wait()

	timerStop <- true
	close(timerStop)

	e.Dest.Done()

	return true
}

func (e *Engine) subRender(origin *common.Vector, startX, stopX, startY, stopY int,
	wg *sync.WaitGroup) {
	defer wg.Done()

	SY := e.WY1 + e.DiffY*float64(startY)

	for y := startY; y <= stopY; y++ {
		SX := e.WX1 + e.DiffX*float64(startX)
		for x := startX; x <= stopX; x++ {

			dir := common.NewVector(SX, SY, 0).Minus(origin)
			dir.Normalize()

			r := common.NewRay(*origin, *dir)

			// if x == 290 && y == 624 {
			// 	r.Debug = true
			// }

			_, _, accColor := e.Raytrace(r, 1)

			e.Dest.Set(x, y, accColor)

			SX += e.DiffX

		}

		SY += e.DiffY
	}
}

func NewEngine() *Engine {
	eng := new(Engine)
	eng.Scene = scene.NewScene()
	return eng
}
