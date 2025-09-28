package scene

import (
	"fmt"

	"github.com/ironsmile/raytracer/accel"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/primitive"
	"github.com/ironsmile/raytracer/scene/example"
)

// Scene is a type which is responsible for loading and managing a scene for rendering.
type Scene struct {
	Primitives []primitive.Primitive
	Lights     []primitive.Primitive
	accel      primitive.Primitive
}

// GetNrLights returns the number of lights in this scene
func (s *Scene) GetNrLights() int {
	return len(s.Lights)
}

// GetLight returns the nth light in the scene
func (s *Scene) GetLight(n int) primitive.Primitive {
	return s.Lights[n]
}

// GetNrPrimitives returns the number of primitives in the scene. This number includes
// all the lights.
func (s *Scene) GetNrPrimitives() int {
	return len(s.Primitives)
}

// GetPrimitive returns the nth primitive in the scene
func (s *Scene) GetPrimitive(n int) primitive.Primitive {
	return s.Primitives[n]
}

// Intersect intersects a ray against all the primitives in the scene.
func (s *Scene) Intersect(ray geometry.Ray, in *primitive.Intersection) bool {
	return s.accel.Intersect(ray, in)
}

// IntersectP tells whether a ray intersects *any* of the primitives in the scene. Thus
// it is faster than `Intersect`.
func (s *Scene) IntersectP(ray geometry.Ray) bool {
	return s.accel.IntersectP(ray)
}

// IntersectBBoxEdge tells whether a ray intersects a bounding box edge of any of the
// primitives in the scene.
func (s *Scene) IntersectBBoxEdge(ray geometry.Ray) bool {
	for _, pr := range s.Primitives {
		if pr.IntersectBBoxEdge(ray) {
			return true
		}
	}
	return false
}

// InitScene programatically creates and loads a demo scene
func (s *Scene) InitScene(name string) {
	var (
		prims  []primitive.Primitive
		lights []primitive.Primitive
	)

	switch name {
	case "car":
		prims, lights = example.GetCarScene()
	case "empty":
		prims, lights = example.GetEmptyScene()
	default:
		prims, lights = example.GetTeapotScene()
	}

	s.Lights = lights
	s.Primitives = prims

	if err := s.initDebugRays(); err != nil {
		fmt.Printf("error adding debug rays to the scene: %s\n", err)
	}

	s.accel = accel.NewBVH(s.Primitives, 1)
	// s.accel = accel.NewGrid(s.Primitives)
}

// NewScene returns a new demo scene
func NewScene() *Scene {
	scn := new(Scene)
	return scn
}

// PossibleScenes is a list of scenes which are supported by `InitScene`.
var PossibleScenes = []string{
	"car",
	"teapot",
	"empty",
}
