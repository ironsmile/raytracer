package engine

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ironsmile/raytracer/camera"
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/sampler"
	"github.com/ironsmile/raytracer/scene"
)

const (
	EPSION     = 0.00001
	TRACEDEPTH = 9
)

type Engine struct {
	Scene         *scene.Scene
	Dest          film.Film
	Width, Height int
	Camera        camera.Camera
	Sampler       sampler.Sampler
}

func (e *Engine) SetTarget(target film.Film, cam camera.Camera) {
	e.Width = target.Width()
	e.Height = target.Height()
	e.Dest = target
	e.Camera = cam
}

func (e *Engine) InitRender(smpl sampler.Sampler) {
	e.Scene = scene.NewScene()
	e.Sampler = smpl
	fmt.Printf("Engine initialized with viewport %dx%d\n", e.Width, e.Height)
}

func (e *Engine) Raytrace(ray *geometry.Ray, depth int64, retColor *geometry.Color) (
	scene.Primitive, float64, *geometry.Color) {

	retColor.Set(0, 0, 0)

	if depth > TRACEDEPTH {
		return nil, 0, retColor
	}

	prim, retdist := e.Scene.Intersect(ray)

	if prim == nil {
		return nil, 0, retColor
	}

	if prim.IsLight() {
		clr := prim.GetColor()
		retColor.Set(clr.Red(), clr.Green(), clr.Blue())
		return prim, retdist, retColor
	}

	pi := ray.Origin.PlusVector(ray.Direction.MultiplyScalar(retdist))

	// if ray.Debug {
	// 	fmt.Printf("I did hit %s\n", prim.GetName())
	// 	ray.Debug = false
	// }

	primMat := prim.GetMaterial()

	for l := 0; l < e.Scene.GetNrLights(); l++ {
		N := prim.GetNormal(pi)
		light := e.Scene.GetLight(l)
		shade := 1.0

		L := (light.(*scene.Sphere)).Center.Minus(pi).NormalizeIP()

		if light.GetType() == scene.SPHERE {
			piOffset := pi.PlusVector(L.MultiplyScalar(EPSION))

			shadowRay := &geometry.Ray{Origin: piOffset, Direction: L}
			// shadowRay.Debug = ray.Debug

			intersected, _ := e.Scene.Intersect(shadowRay)

			if light != intersected {
				shade = 0.0
			}
		}

		if primMat.Diff > 0 {
			dot := N.Product(L)
			if dot > 0 {
				weight := dot * primMat.Diff * shade
				retColor.PlusIP(light.GetMaterial().Color.
					Multiply(primMat.Color).MultiplyScalarIP(weight))
			}
		}

		if primMat.GetSpecular() > 0 {
			V := ray.Direction
			R := L.Minus(N.MultiplyScalar(2.0 * L.Product(N)))
			dot := V.Product(R)
			if dot > 0 {
				spec := math.Pow(dot, 20) * primMat.GetSpecular() * shade
				retColor.PlusIP(light.GetMaterial().Color.
					MultiplyScalar(spec))
			}
		}

	}

	// Reflection
	if primMat.Refl > 0.0 {
		N := prim.GetNormal(pi)
		R := ray.Direction.Minus(N.MultiplyScalarIP(ray.Direction.Product(N) * 2.0))

		refRay := &geometry.Ray{Origin: pi.PlusVector(R.MultiplyScalar(EPSION)),
			Direction: R}
		// refRay.Debug = ray.Debug
		refColor := &geometry.Color{}
		e.Raytrace(refRay, depth+1, refColor)

		retColor.PlusIP(primMat.Color.Multiply(
			refColor).MultiplyScalarIP(primMat.Refl))
	}

	return prim, retdist, retColor
}

func (e *Engine) startParallelRendering(wg *sync.WaitGroup,
	subRendererFunction func(int, int, int, int, *sync.WaitGroup)) {
	quads := 3
	quadWidth := e.Width / quads
	quadHeight := e.Height / quads

	for quadIndX := 0; quadIndX < quads; quadIndX++ {
		for quadIndY := 0; quadIndY < quads; quadIndY++ {

			quadXStart := quadIndX * quadWidth
			quadXStop := quadXStart + quadWidth - 1

			quadYStart := quadIndY * quadHeight
			quadYStop := quadYStart + quadHeight - 1

			wg.Add(1)
			go subRendererFunction(quadXStart, quadXStop, quadYStart, quadYStop, wg)
		}
	}
}

func (e *Engine) Render() {

	var wg sync.WaitGroup

	e.Dest.StartFrame()

	engineTimer := time.Now()

	for i := 0; i < 9; i++ {
		wg.Add(1)
		go e.subRender(&wg)
	}

	wg.Wait()

	fmt.Printf("Engine frame time: %s\n", time.Since(engineTimer))

	e.Dest.DoneFrame()
}

func (e *Engine) subRender(wg *sync.WaitGroup) {
	defer wg.Done()
	ray := &geometry.Ray{}
	accColor := geometry.NewColor(0, 0, 0)

	for {
		x, y, err := e.Sampler.GetSample()
		if err != nil {
			return
		}
		weight := e.Camera.GenerateRayIP(float64(x), float64(y), ray)
		e.Raytrace(ray, 1, accColor)
		e.Sampler.UpdateScreen(x, y, accColor.MultiplyScalarIP(weight))
	}

}

func NewEngine() *Engine {
	eng := new(Engine)
	return eng
}
