package engine

import (
	"fmt"
	"math/rand"
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
	Sampler       *sampler.SimpleSampler
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
func (e *Engine) Raytrace(ray geometry.Ray, depth int64, in *primitive.Intersection) geometry.Color {
	var retColor geometry.Color

	if depth > TraceDepth {
		return retColor
	}

	if ok := e.Scene.Intersect(ray, in); !ok {
		return retColor
	}

	prim := in.Primitive
	pi := ray.At(in.DfGeometry.Distance)

	if prim.IsLight() {
		return *prim.Shape().MaterialAt(pi).Color
	}

	o2w, w2o := prim.GetTransforms()
	pio := w2o.Point(pi)
	inNormal := o2w.Normal(in.DfGeometry.Shape.NormalAt(pio))

	cosI := inNormal.Dot(ray.Direction)
	if cosI > 0 {
		// The hit is from the inside of the primitive. Normally, all normals would be
		// pointing toward the primitive exterior. So we have to invert it to the interior
		// for proper calculations.
		inNormal = inNormal.Neg()
	}

	primMat := in.DfGeometry.Shape.MaterialAt(pio)

	// /* Debugging */
	// var debugging bool
	// if !e.debugged && ray.Debug {
	// 	e.debugged = true
	// 	debugging = true
	// 	fmt.Printf("\nIntersected: %s\nnormal: %s\nretdist: %f\n",
	// 		prim.GetName(), inNormal, retdist)
	// }

	light := e.calculateLight(pi, inNormal)

	if primMat.Diff > 0 {
		retColor.PlusIP(primMat.Color.Multiply(&light))
	}

	// Reflection
	if primMat.Refl > 0.0 {

		R := ray.Direction.Minus(inNormal.MultiplyScalar(
			ray.Direction.Product(inNormal) * 2.0),
		)

		refRay := geometry.NewRay(pi, R)
		refRay.Mint = geometry.EPSILON

		// refRay.Debug = ray.Debug
		refColor := e.Raytrace(refRay, depth+1, in)

		retColor.PlusIP(primMat.Color.Multiply(
			&refColor).MultiplyScalarIP(primMat.Refl))
	}

	// Refraction
	if primMat.Refr > 0.0 && primMat.RefrIndex > 0 {

		var refrNormal = inNormal
		var n1, n2 = 1.0, primMat.RefrIndex
		var reflectance, transmittance float64

		refrDirection, tir := ray.Refract(refrNormal, n1, n2)

		if tir {
			reflectance = 1
		} else {
			reflectance = geometry.Schlick2(refrNormal, ray.Direction, n1, n2)
		}

		var endColor geometry.Color

		transmittance = 1 - reflectance

		if transmittance > 0 {
			reflRay := geometry.NewRay(pi, refrDirection)
			reflRay.Mint = geometry.EPSILON * 2
			refrColor := e.Raytrace(reflRay, depth+1, in)
			endColor.PlusIP(refrColor.MultiplyScalarIP(transmittance))
		}

		if reflectance > 0 {
			cosI := -refrNormal.Dot(ray.Direction)
			R := ray.Direction.Plus(refrNormal.MultiplyScalar(2 * cosI))

			refRay := geometry.NewRay(pi, R)
			refRay.Mint = geometry.EPSILON
			refColor := e.Raytrace(refRay, depth+1, in)
			endColor.PlusIP(refColor.MultiplyScalarIP(reflectance))
		}

		retColor.MultiplyScalarIP(1 - primMat.Refr).PlusIP(
			(&endColor).MultiplyScalarIP(primMat.Refr))
	}

	return retColor
}

func (e *Engine) calculateLight(pi, inNormal geometry.Vector) geometry.Color {
	l := rand.Intn(e.Scene.GetNrLights() - 1)
	light := e.Scene.GetLight(l)

	source := light.GetLightSource()
	shadowRayStart := pi.Plus(inNormal.MultiplyScalar(geometry.EPSILON))
	L := source.Minus(shadowRayStart).Normalize()
	shadowRay := geometry.NewRay(shadowRayStart, L)
	shadowRay.Maxt = shadowRayStart.Distance(source)

	if intersected := e.Scene.IntersectP(shadowRay); intersected {
		return geometry.Black
	}

	dot := inNormal.Product(L)

	if dot <= 0 {
		return geometry.Black
	}

	return *light.Shape().MaterialAt(source).Color.MultiplyScalar(dot)
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

	var accColor geometry.Color
	var in primitive.Intersection

	for {

		subSampler, err := e.Sampler.GetSubSampler()

		if err == sampler.ErrEndOfSampling {
			return
		}

		for {
			x, y, err := subSampler.GetSample()

			if err == sampler.ErrEndOfSampling {
				return
			}
			if err == sampler.ErrSubSamplerEnd {
				break
			}
			if err != nil {
				fmt.Printf("Error while getting sample: %s\n", err)
				return
			}

			// fmt.Printf("x: %f, y: %f\n", x, y)

			ray := e.Camera.GenerateRay(x, y)
			accColor = e.Raytrace(ray, 1, &in)

			if e.ShowBBoxes {
				if in.Primitive != nil {
					ray.Maxt = in.DfGeometry.Distance
				}

				if e.Scene.IntersectBBoxEdge(ray) {
					accColor = *geometry.NewColor(0, 0, 1)
				}

				// debugRay := geometry.NewRay(
				// 	geometry.NewVector(-10.000000, -0.097124, 0.562618),
				// 	geometry.NewVector(0.151860, -0.768368, 0.621731),
				// )
				// if _, ok := ray.Intersect(debugRay); ok {
				// 	accColor = *geometry.NewColor(1, 1, 0)
				// }
			}

			e.Sampler.UpdateScreen(x, y, &accColor)
		}
	}
}

// New returns a new engine which would use the argument's sampler
func New(smpl *sampler.SimpleSampler) *Engine {
	eng := new(Engine)
	initEngine(eng, smpl)
	return eng
}

func initEngine(eng *Engine, smpl *sampler.SimpleSampler) {
	eng.Scene = scene.NewScene()
	eng.Sampler = smpl
}
