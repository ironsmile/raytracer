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
	EPSION     = 0.00001
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

	retColor.Set(0, 0, 0)

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
	L := &geometry.Vector{}
	piOffset := &geometry.Point{}

	piDirection := ray.Direction.MultiplyScalar(retdist)
	pi := ray.Origin.PlusVector(piDirection)

	primMat := prim.GetMaterial()

	/* Debugging */
	// if !e.debugged && prim.GetName() == "First teapod" {
	// 	fmt.Printf("clr %s: %v\n", prim.GetName(), primMat.Color)
	// 	fmt.Printf("InNormal %s\n", InNormal)
	// 	ray.Debug = true
	// 	e.debugged = true
	// }

	//!TODO: maybe make sure the Intersect method returns a copy of a internal Normal vector
	// so that it can be modified. InNormal, for example is modified further down this method.
	// InNormal = InNormal.Copy()

	for l := 0; l < e.Scene.GetNrLights(); l++ {
		light := e.Scene.GetLight(l)
		luminousity := 1.0

		// Reusing the same object as much as possible
		light.GetLightSource().MinusInVector(pi, L)
		L.NormalizeIP()

		dot := InNormal.Product(L)

		if light.GetType() == primitive.SPHERE {
			piOffset.CopyToSelf(pi).PlusVectorIP(L.MultiplyScalar(EPSION))

			// Reusing the same object as much as possible
			shadowRay.BackToDefaults()
			shadowRay.Origin = *piOffset
			shadowRay.Direction = *L

			// shadowRay.Debug = ray.Debug

			intersected, _, _ := e.Scene.Intersect(shadowRay)

			if light != intersected {
				// This was previously `luminousity = 0` which gave a really hard
				// shadowing. In an effort to make the lightning softer I changed
				// it to the one ti is now. I don't really know how this will
				// affect any other parts of the tracer.
				luminousity = 1.0 - dot
			}
		}

		if primMat.Diff > 0 && dot > 0 {
			weight := dot * primMat.Diff * luminousity
			retColor.PlusIP(light.GetMaterial().Color.
				Multiply(primMat.Color).MultiplyScalarIP(weight))
		}

		if primMat.GetSpecular() > 0 {
			V := ray.Direction
			R := L.Minus((&InNormal).MultiplyScalar(2.0 * L.Product(&InNormal)))
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

		// Warning! InNormal is irrevrsibly changed here.
		R := ray.Direction.Minus((&InNormal).MultiplyScalarIP(ray.Direction.Product(&InNormal) * 2.0))

		refRay := geometry.Ray{
			Origin:    *pi.PlusVectorIP(R.MultiplyScalar(EPSION)),
			Direction: *R,
		}

		// refRay.Debug = ray.Debug
		_, _, refColor := e.Raytrace(refRay, depth+1)

		retColor.PlusIP(primMat.Color.Multiply(
			&refColor).MultiplyScalarIP(primMat.Refl))
	}

	if e.ShowBBoxes && depth == 1 && e.Scene.IntersectBBoxEdge(ray, retdist) {
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
	ray := geometry.Ray{}
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
		weight := e.Camera.GenerateRayIP(float64(x), float64(y), &ray)
		_, _, accColor = e.Raytrace(ray, 1)
		e.Sampler.UpdateScreen(x, y, accColor.MultiplyScalarIP(weight))
	}

}

func New(smpl sampler.Sampler) *Engine {
	eng := new(Engine)
	initEngine(eng, smpl)
	return eng
}

func initEngine(eng *Engine, smpl sampler.Sampler) {
	eng.Scene = scene.NewScene()
	eng.Sampler = smpl
}
