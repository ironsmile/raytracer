package scene

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
)

const (
	HIT = iota
	MISS
	INPRIM
)

const (
	NOTHING = iota
	SPHERE
	PLANE
)

type Scene struct {
	Primitives []Primitive
	Lights     []Primitive
}

func (s *Scene) GetNrLights() int {
	return len(s.Lights)
}

func (s *Scene) GetLight(index int) Primitive {
	return s.Lights[index]
}

func (s *Scene) GetNrPrimitives() int {
	return len(s.Primitives)
}

func (s *Scene) GetPrimitive(index int) Primitive {
	return s.Primitives[index]
}

func (s *Scene) Intersect(ray *geometry.Ray) (Primitive, float64) {
	retdist := 1000000.0
	var prim Primitive = nil

	for sInd := 0; sInd < len(s.Primitives); sInd++ {
		pr := s.Primitives[sInd]

		if pr == nil {
			fmt.Errorf("Primitive with index %d was nil\n", sInd)
		}

		res, resDist := pr.Intersect(ray, retdist)

		if res != MISS {
			prim = pr
			retdist = resDist
		}
	}

	return prim, retdist
}

func (s *Scene) InitScene() {
	s.Primitives = make([]Primitive, 0)
	s.Lights = make([]Primitive, 0)

	plane := NewPlanePrim(*geometry.NewVector(0, 1, 0), 4)
	plane.Name = "plane"
	plane.Mat.Refl = 0
	plane.Mat.Diff = 1.0
	plane.Mat.Color = geometry.NewColor(0.4, 0.3, 0.3)

	s.Primitives = append(s.Primitives, plane)

	plane = NewPlanePrim(*geometry.NewVector(1, 0, 0), 11)
	plane.Name = "plane"
	plane.Mat.Refl = 0
	plane.Mat.Diff = 1.0
	plane.Mat.Color = geometry.NewColor(0.4, 0.3, 0.3)

	s.Primitives = append(s.Primitives, plane)

	sphere := NewSphere(*geometry.NewPoint(1, -0.8, 3), 2.5)
	sphere.Name = "big sphere"
	sphere.Mat.Refl = 0.8
	sphere.Mat.Diff = 0.9
	sphere.Mat.Color = geometry.NewColor(1, 0, 0)

	s.Primitives = append(s.Primitives, sphere)

	sphere = NewSphere(*geometry.NewPoint(-5.5, -0.5, 7), 2)
	sphere.Name = "small sphere"
	sphere.Mat.Refl = 0.9
	sphere.Mat.Diff = 0.4
	sphere.Mat.Color = geometry.NewColor(0.7, 0.7, 1)

	s.Primitives = append(s.Primitives, sphere)

	sphere = NewSphere(*geometry.NewPoint(-6.5, -2.5, 25), 1.5)
	sphere.Name = "small sphere far away"
	sphere.Mat.Refl = 0.9
	sphere.Mat.Diff = 0.4
	sphere.Mat.Color = geometry.NewColor(0.5, 1, 0)

	s.Primitives = append(s.Primitives, sphere)

	sphere = NewSphere(*geometry.NewPoint(0, 5, 5), 0.1)
	sphere.Name = "Visible light source"
	sphere.Light = true
	sphere.Mat.Color = geometry.NewColor(0.9, 0.9, 0.9)

	s.Primitives = append(s.Primitives, sphere)
	s.Lights = append(s.Lights, sphere)

	sphere = NewSphere(*geometry.NewPoint(2, 5, 1), 0.1)
	sphere.Name = "Invisible lightsource"
	sphere.Light = true
	sphere.Mat.Color = geometry.NewColor(0.9, 0.9, 0.9)

	s.Primitives = append(s.Primitives, sphere)
	s.Lights = append(s.Lights, sphere)
}

func NewScene() *Scene {
	scn := new(Scene)
	return scn
}
