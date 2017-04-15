package engine

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/ironsmile/raytracer/camera"
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/primitive"
	"github.com/ironsmile/raytracer/sampler"
	"github.com/ironsmile/raytracer/scene"
)

const (
	TRACEDEPTH = 9
)

type Engine struct {
	Scene         *scene.Scene
	Dest          film.Film
	Width, Height int
	Camera        camera.Camera
	Sampler       sampler.Sampler
	ShowBBoxes    bool

	debugged bool
}

func (e *Engine) SetTarget(target film.Film, cam camera.Camera) {
	e.Width = target.Width()
	e.Height = target.Height()
	e.Dest = target
	e.Camera = cam
}

func (e *Engine) Raytrace(ray geometry.Ray, depth int64) (
	primitive.Primitive, float64, geometry.Color) {
	var retColor geometry.Color

	if depth > TRACEDEPTH {
		return nil, 0, retColor
	}

	prim, retdist, InNormal := e.Scene.Intersect(ray)

	if prim == nil {
		return nil, 0, retColor
	}

	if prim.IsLight() {
		clr := prim.GetColor()
		retColor.Set(clr.Red(), clr.Green(), clr.Blue())
		return prim, retdist, retColor
	}

	shadowRay := geometry.Ray{}

	pi := ray.At(retdist)

	primMat := prim.GetMaterial()

	// /* Debugging */
	// var debugging bool
	// if !e.debugged && ray.Debug {
	// 	e.debugged = true
	// 	debugging = true
	// 	fmt.Printf("\nIntersected: %s\nnormal: %s\nretdist: %f\n",
	// 		prim.GetName(), InNormal, retdist)
	// }

	for l := 0; l < e.Scene.GetNrLights(); l++ {
		light := e.Scene.GetLight(l)
		luminousity := 0.0

		// Reusing the same object as much as possible
		L := light.GetLightSource().Minus(pi).Normalize()

		dot := InNormal.Product(L)

		if light.GetType() == primitive.SPHERE {

			// Reusing the same object as much as possible
			shadowRay.BackToDefaults()
			shadowRay.Origin = pi
			shadowRay.Direction = L
			shadowRay.Maxt = pi.Distance(light.GetLightSource())
			shadowRay.Mint = geometry.EPSILON

			intersected, _, _ := e.Scene.Intersect(shadowRay)

			if light == intersected {
				luminousity = 0.8
			}
		}

		if luminousity > 0 && primMat.Diff > 0 && dot > 0 {
			weight := dot * primMat.Diff * luminousity
			retColor.PlusIP(light.GetMaterial().Color.
				Multiply(primMat.Color).MultiplyScalarIP(weight))
		}

		if luminousity > 0 && primMat.GetSpecular() > 0 {
			V := ray.Direction
			R := L.Minus(InNormal.MultiplyScalar(2.0 * L.Product(InNormal)))
			dot := V.Product(R)
			if dot > 0 {
				spec := math.Pow(dot, 20) * primMat.GetSpecular() * luminousity
				retColor.PlusIP(light.GetMaterial().Color.
					MultiplyScalar(spec))
			}
		}
	}

	// Reflection
	if primMat.Refl > 0.0 {

		R := ray.Direction.Minus(InNormal.MultiplyScalar(ray.Direction.Product(InNormal) * 2.0))

		refRay := geometry.NewRay(pi, R)
		refRay.Mint = geometry.EPSILON

		// refRay.Debug = ray.Debug
		_, _, refColor := e.Raytrace(refRay, depth+1)

		retColor.PlusIP(primMat.Color.Multiply(
			&refColor).MultiplyScalarIP(primMat.Refl))
	}

	ray.Maxt = retdist
	if e.ShowBBoxes && depth == 1 && e.Scene.IntersectBBoxEdge(ray) {
		retColor = *geometry.NewColor(0, 0, 1)
	}

	return prim, retdist, retColor
}

func (e *Engine) Render() {

	var wg sync.WaitGroup

	e.Dest.StartFrame()

	engineTimer := time.Now()

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go e.subRender(&wg)
	}

	wg.Wait()

	fmt.Printf("Engine frame time: %s\n", time.Since(engineTimer))

	e.Dest.DoneFrame()
}

func (e *Engine) subRender(wg *sync.WaitGroup) {
	defer wg.Done()

	var accColor geometry.Color

	for {
		x, y, err := e.Sampler.GetSample()
		if err == sampler.EndOfSampling {
			return
		}
		if err != nil {
			fmt.Printf("Error while getting sample: %s\n", err)
			return
		}
		ray, weight := e.Camera.GenerateRay(x, y)
		_, _, accColor = e.Raytrace(ray, 1)
		e.Sampler.UpdateScreen(x, y, accColor.MultiplyScalarIP(weight))
	}
}

// New returns a new engine which would use the argument's sampler
func New(smpl sampler.Sampler) *Engine {
	eng := new(Engine)
	initEngine(eng, smpl)
	return eng
}

func initEngine(eng *Engine, smpl sampler.Sampler) {
	eng.Scene = scene.NewScene()
	eng.Sampler = smpl
}
