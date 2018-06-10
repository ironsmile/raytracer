package engine

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/ironsmile/raytracer/camera"
	"github.com/ironsmile/raytracer/color"
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"
	"github.com/ironsmile/raytracer/primitive"
	"github.com/ironsmile/raytracer/sampler"
	"github.com/ironsmile/raytracer/scene"
)

const (
	// TraceDepth is the limit of generated rays recursion
	TraceDepth = 5
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
func (e *Engine) Raytrace(
	ray geometry.Ray,
	depth int64,
	rnd *rand.Rand,
	in *primitive.Intersection) color.Color {

	if depth > TraceDepth {
		return color.Black
	}

	if ok := e.Scene.Intersect(ray, in); !ok {
		return color.Sky
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

	if true {
		directLight := e.calculateLight(pi, inNormal, rnd)

		indirectRay := e.getBRDFRay(ray, primMat, inNormal, pi, rnd)
		indirectLight := e.Raytrace(indirectRay, depth+1, rnd, in)

		return *primMat.Color.Multiply(directLight.Plus(&indirectLight))
	}

	indirectRay := e.getBRDFRay(ray, primMat, inNormal, pi, rnd)
	indirectLight := e.Raytrace(indirectRay, depth+1, rnd, in)

	return *primMat.Color.Multiply(&indirectLight)
}

func (e *Engine) getBRDFRay(
	ray geometry.Ray,
	primMat *mat.Material,
	inNormal geometry.Vector,
	pi geometry.Vector,
	rnd *rand.Rand,
) geometry.Ray {

	chance := rnd.Float64()

	cosI := -inNormal.Dot(ray.Direction)
	reflectionDirection := ray.Direction.Plus(inNormal.MultiplyScalar(2 * cosI))

	// Reflection
	if chance <= primMat.Refl {
		refRay := geometry.NewRay(pi, reflectionDirection)
		refRay.Mint = geometry.EPSILON
		return refRay
	}

	// Refraction
	if primMat.Refr > 0.0 && primMat.RefrIndex > 0 {
		var n1, n2 = 1.0, primMat.RefrIndex
		var reflectance float64

		refrDirection, tir := ray.Refract(inNormal, n1, n2)

		if tir {
			reflectance = 1
		} else {
			reflectance = geometry.Schlick2(inNormal, ray.Direction, n1, n2)
		}

		chance = rnd.Float64()

		if chance <= reflectance {
			refRay := geometry.NewRay(pi, reflectionDirection)
			refRay.Mint = geometry.EPSILON
			return refRay
		}

		// transmittance
		reflRay := geometry.NewRay(pi, refrDirection)
		reflRay.Mint = geometry.EPSILON * 2
		return reflRay
	}

	// Generating a ray in a hemiosphere around the normal of intersection
	rndDirection := geometry.Vector{
		X: rnd.Float64() - 0.5,
		Y: rnd.Float64() - 0.5,
		Z: rnd.Float64() - 0.5,
	}

	if rndDirection.Product(inNormal) < 0 {
		rndDirection = rndDirection.Neg()
	}

	refRay := geometry.NewRay(pi, rndDirection)
	refRay.Mint = geometry.EPSILON

	return refRay
}

func (e *Engine) calculateLight(pi, inNormal geometry.Vector, rnd *rand.Rand) color.Color {
	l := rnd.Intn(e.Scene.GetNrLights() - 1)
	light := e.Scene.GetLight(l)
	in := primitive.Intersection{}

	source := light.GetLightSource()
	shadowRayStart := pi.Plus(inNormal.MultiplyScalar(geometry.EPSILON))
	L := source.Minus(shadowRayStart).Normalize()
	shadowRay := geometry.NewRay(shadowRayStart, L)
	shadowRay.Maxt = shadowRayStart.Distance(source)

	if intersected := e.Scene.Intersect(shadowRay, &in); !intersected {
		return color.Black
	}

	if in.Primitive != light {
		return color.Black
	}

	dot := inNormal.Product(L)

	if dot <= 0 {
		return color.Black
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

	var accColor color.Color
	var in primitive.Intersection

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {

		subSampler, err := e.Sampler.GetSubSampler()

		if err == sampler.ErrEndOfSampling {
			return
		}

		for {
			x, y, err := subSampler.GetSample(rnd)

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
			accColor = e.Raytrace(ray, 1, rnd, &in)

			if e.ShowBBoxes {
				if in.Primitive != nil {
					ray.Maxt = in.DfGeometry.Distance
				}

				if e.Scene.IntersectBBoxEdge(ray) {
					accColor = *color.NewColor(0, 0, 1)
				}

				// debugRay := geometry.NewRay(
				// 	geometry.NewVector(-10.000000, -0.097124, 0.562618),
				// 	geometry.NewVector(0.151860, -0.768368, 0.621731),
				// )
				// if _, ok := ray.Intersect(debugRay); ok {
				// 	accColor = *color.NewColor(1, 1, 0)
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
