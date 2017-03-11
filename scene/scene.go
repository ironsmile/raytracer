package scene

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/primitive"
	"github.com/ironsmile/raytracer/shape"
)

type Scene struct {
	Primitives []primitive.Primitive
	Lights     []primitive.Primitive
}

func (s *Scene) GetNrLights() int {
	return len(s.Lights)
}

func (s *Scene) GetLight(index int) primitive.Primitive {
	return s.Lights[index]
}

func (s *Scene) GetNrPrimitives() int {
	return len(s.Primitives)
}

func (s *Scene) GetPrimitive(index int) primitive.Primitive {
	return s.Primitives[index]
}

func (s *Scene) Intersect(ray *geometry.Ray) (prim primitive.Primitive, retdist float64, normal *geometry.Vector) {
	retdist = 1000000.0

	for sInd, pr := range s.Primitives {

		if pr == nil {
			fmt.Printf("primitive with index %d was nil\n", sInd)
			continue
		}

		res, resDist, resNormal := pr.Intersect(ray, retdist)

		if res == shape.HIT && resDist < retdist {
			prim = pr
			retdist = resDist
			normal = resNormal
		}
	}

	return
}

func (s *Scene) InitScene() {
	s.Primitives = make([]primitive.Primitive, 0)
	s.Lights = make([]primitive.Primitive, 0)

	plane := primitive.NewPlanePrim(geometry.NewVector(0, 1, 0), 4)
	plane.Name = "plane-floor"
	plane.Mat.Refl = 0
	plane.Mat.Diff = 1.0
	plane.Mat.Color = geometry.NewColor(0.4, 0.3, 0.3)

	s.Primitives = append(s.Primitives, plane)

	plane = primitive.NewPlanePrim(geometry.NewVector(0, -1, 0), 11)
	plane.Name = "plane-ceiling"
	plane.Mat.Refl = 0
	plane.Mat.Diff = 1.0
	plane.Mat.Color = geometry.NewColor(0.4, 0.3, 0.3)

	s.Primitives = append(s.Primitives, plane)

	plane = primitive.NewPlanePrim(geometry.NewVector(1, 0, 0), 33)
	plane.Name = "plane-left"
	plane.Mat.Refl = 0
	plane.Mat.Diff = 1.0
	plane.Mat.Color = geometry.NewColor(0.4, 0.3, 0.3)

	s.Primitives = append(s.Primitives, plane)

	plane = primitive.NewPlanePrim(geometry.NewVector(-1, 0, 0), 11)
	plane.Name = "plane-right-mirror"
	plane.Mat.Refl = 1.0
	plane.Mat.Diff = 0.4
	plane.Mat.Color = geometry.NewColor(0.4, 0.3, 0.3)

	s.Primitives = append(s.Primitives, plane)

	plane = primitive.NewPlanePrim(geometry.NewVector(0, 0, -1), 30)
	plane.Name = "plane-front"
	plane.Mat.Refl = 0
	plane.Mat.Diff = 1.0
	plane.Mat.Color = geometry.NewColor(0.4, 0.3, 0.3)

	s.Primitives = append(s.Primitives, plane)

	sphere := primitive.NewSphere(*geometry.NewPoint(1, -0.8, 3), 2.5)
	sphere.Name = "big sphere"
	sphere.Mat.Refl = 0.8
	sphere.Mat.Diff = 0.9
	sphere.Mat.Color = geometry.NewColor(1, 0, 0)

	s.Primitives = append(s.Primitives, sphere)

	sphere = primitive.NewSphere(*geometry.NewPoint(-5.5, -0.5, 7), 2)
	sphere.Name = "small sphere"
	sphere.Mat.Refl = 0.9
	sphere.Mat.Diff = 0.4
	sphere.Mat.Color = geometry.NewColor(0.7, 0.7, 1)

	s.Primitives = append(s.Primitives, sphere)

	sphere = primitive.NewSphere(*geometry.NewPoint(-6.5, -2.5, 25), 1.5)
	sphere.Name = "small sphere far away"
	sphere.Mat.Refl = 0.9
	sphere.Mat.Diff = 0.4
	sphere.Mat.Color = geometry.NewColor(0.5, 1, 0)

	s.Primitives = append(s.Primitives, sphere)

	triangle := primitive.NewTriangle([3]*geometry.Point{
		geometry.NewPoint(-10.99, 3, 0),  // a
		geometry.NewPoint(-10.99, 0, -3), // b
		geometry.NewPoint(-10.99, 0, 3),  // c
	})
	triangle.Name = "Green triangle"
	triangle.Mat.Refl = 0.0
	triangle.Mat.Diff = 0.3
	triangle.Mat.Color = geometry.NewColor(0.3, 1, 0)

	s.Primitives = append(s.Primitives, triangle)

	sphere = primitive.NewSphere(*geometry.NewPoint(0, 5, 5), 0.1)
	sphere.Name = "Visible light source"
	sphere.Light = true
	sphere.Mat.Color = geometry.NewColor(0.9, 0.9, 0.9)

	s.Primitives = append(s.Primitives, sphere)
	s.Lights = append(s.Lights, sphere)

	sphere = primitive.NewSphere(*geometry.NewPoint(2, 5, 1), 0.1)
	sphere.Name = "Invisible lightsource"
	sphere.Light = true
	sphere.Mat.Color = geometry.NewColor(0.9, 0.9, 0.9)

	s.Primitives = append(s.Primitives, sphere)
	s.Lights = append(s.Lights, sphere)

	teaPotCenter := geometry.NewPoint(-3, 0, 5)
	if teapot, err := primitive.NewObject("data/objs/teapot.obj", teaPotCenter); err != nil {
		fmt.Printf("Error loading obj teapot: %s\n", err)
	} else {
		teapot.Name = "First teapod"
		teapot.Mat.Refl = 0.0
		teapot.Mat.Diff = 0.3
		teapot.Mat.Color = geometry.NewColor(0.3, 1, 0)
		s.Primitives = append(s.Primitives, teapot)
	}
}

func NewScene() *Scene {
	scn := new(Scene)
	return scn
}
