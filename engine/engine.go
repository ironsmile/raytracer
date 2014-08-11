package engine

import (
	"fmt"
	"math"
	"sync"

	"github.com/ironsmile/raytracer/camera"
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/geometry"
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
	Camera                                   camera.Camera
}

func (e *Engine) SetTarget(target film.Film, cam camera.Camera) {
	e.Width = int64(target.Width())
	e.Height = int64(target.Height())
	e.Dest = target
	e.Camera = cam
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

func (e *Engine) Raytrace(ray *geometry.Ray, depth int64) (
	scene.Primitive, float64, *geometry.Color) {

	retColor := geometry.NewColor(0, 0, 0)

	if depth > geometry.TRACEDEPTH {
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

	pi := ray.Origin.PlusVector(ray.Direction.MultiplyScalar(retdist))

	for l := 0; l < e.Scene.GetNrLights(); l++ {
		light := e.Scene.GetLight(l)
		shade := 1.0

		L := (light.(*scene.Sphere)).Center.Minus(pi)
		L.Normalize()

		if light.GetType() == scene.SPHERE {
			piOffset := pi.PlusVector(L.MultiplyScalar(EPSION))
			shadowRay := &geometry.Ray{Origin: piOffset, Direction: L}
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
				retColor = retColor.Plus(light.GetMaterial().Color.
					Multiply(primMat.Color).MultiplyScalar(weight))
			}
		}

		if primMat.GetSpecular() > 0 {
			V := ray.Direction
			R := L.Minus(N.MultiplyScalar(2.0 * L.Product(N)))
			dot := V.Product(R)
			if dot > 0 {
				spec := math.Pow(dot, 20) * primMat.GetSpecular() * shade
				retColor = retColor.Plus(light.GetMaterial().Color.
					MultiplyScalar(spec))
			}
		}

	}

	// Reflection
	if primMat.Refl > 0.0 {
		N := prim.GetNormal(pi)
		R := ray.Direction.Minus(N.MultiplyScalar(ray.Direction.Product(N) * 2.0))

		refRay := &geometry.Ray{Origin: pi.PlusVector(R.MultiplyScalar(EPSION)),
			Direction: R}
		// refRay.Debug = ray.Debug
		_, _, refColor := e.Raytrace(refRay, depth+1)

		retColor = retColor.Plus(primMat.Color.Multiply(
			refColor).MultiplyScalar(primMat.Refl))
	}

	return prim, retdist, retColor
}

func (e *Engine) Render() bool {

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
			go e.subRender(quadXStart, quadXStop, quadYStart, quadYStop, &wg)
		}
	}

	wg.Wait()

	e.Dest.Done()

	return true
}

func (e *Engine) subRender(startX, stopX, startY, stopY int,
	wg *sync.WaitGroup) {
	defer wg.Done()

	SY := e.WY1 + e.DiffY*float64(startY)

	for y := startY; y <= stopY; y++ {
		SX := e.WX1 + e.DiffX*float64(startX)
		for x := startX; x <= stopX; x++ {

			// r, weight := e.Camera.GenerateRay(SX, SY)
			r, weight := e.Camera.GenerateRay(float64(x), float64(y))

			if x == camera.DEBUG_X && y == camera.DEBUG_Y {
				fmt.Printf("Final ray:\n%v\n", r)
			}

			_, _, accColor := e.Raytrace(r, 1)

			e.Dest.Set(x, y, accColor.MultiplyScalar(weight))

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
