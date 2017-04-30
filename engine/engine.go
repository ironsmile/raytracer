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
	// TraceDepth is the limit of generated rays recursion
	TraceDepth = 9
)

// Engine is the type which is resposible for bringing the camera, scene and
// everything else together. It generates the rays, intersects them and then
// paints the result in the output film.
type Engine struct {
	Scene         *scene.Scene
	Dest          film.Film
	Width, Height int
	Camera        camera.Camera
	Sampler       sampler.Sampler
	ShowBBoxes    bool

	debugged bool
}

// SetTarget sets the camera and film for rendering.
func (e *Engine) SetTarget(target film.Film, cam camera.Camera) {
	e.Width = target.Width()
	e.Height = target.Height()
	e.Dest = target
	e.Camera = cam
}

// Raytrace returns intersection information for particular ray in the engine's
// scene.
func (e *Engine) Raytrace(ray geometry.Ray, depth int64, in *primitive.Intersection) (
	primitive.Primitive, float64, geometry.Color) {
	var retColor geometry.Color

	if depth > TraceDepth {
		return nil, 0, retColor
	}

	if ok := e.Scene.Intersect(ray, in); !ok {
		return nil, 0, retColor
	}

	prim := in.Primitive
	o2w, _ := prim.GetTransforms()
	retdist := in.DfGeometry.Distance
	InNormal := o2w.Normal(in.DfGeometry.Normal)

	if prim.IsLight() {
		return prim, retdist, *prim.GetColor()
	}

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

		shadowRayStart := pi.Plus(InNormal.MultiplyScalar(geometry.EPSILON))
		L := light.GetLightSource().Minus(shadowRayStart).Normalize()
		shadowRay := geometry.NewRay(shadowRayStart, L)
		shadowRay.Maxt = shadowRayStart.Distance(light.GetLightSource())

		if intersected := e.Scene.IntersectP(shadowRay); intersected {
			continue
		}

		dot := InNormal.Product(L)
		luminousity := 0.8

		if primMat.Diff > 0 && dot > 0 {
			weight := dot * primMat.Diff * luminousity
			retColor.PlusIP(light.GetMaterial().Color.
				Multiply(primMat.Color).MultiplyScalarIP(weight))
		}

		if primMat.GetSpecular() > 0 {
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

		R := ray.Direction.Minus(InNormal.MultiplyScalar(
			ray.Direction.Product(InNormal) * 2.0),
		)

		refRay := geometry.NewRay(pi, R)
		refRay.Mint = geometry.EPSILON

		// refRay.Debug = ray.Debug
		_, _, refColor := e.Raytrace(refRay, depth+1, in)

		retColor.PlusIP(primMat.Color.Multiply(
			&refColor).MultiplyScalarIP(primMat.Refl))
	}

	return prim, retdist, retColor
}

// Render starts the rendering process. Exits when one full frame is done. It does that
// by starting multiple concurrent renderer goroutines.
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

	var pr primitive.Primitive
	var accColor geometry.Color
	var in primitive.Intersection

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
		pr, _, accColor = e.Raytrace(ray, 1, &in)

		if e.ShowBBoxes {
			if pr != nil {
				ray.Maxt = in.DfGeometry.Distance
			}

			if e.Scene.IntersectBBoxEdge(ray) {
				accColor = *geometry.NewColor(0, 0, 1)
			}

			// Debug ray visualization example:
			//
			// debugRay := geometry.NewRay(
			// 	geometry.NewVector(-10.000000, 0.030721, 0.435039),
			// 	geometry.NewVector(0.882750, -0.448091, -0.141302),
			// )

			// if _, ok := ray.Intersect(debugRay); ok {
			// 	accColor = *geometry.NewColor(0, 0, 1)
			// }
		}

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
