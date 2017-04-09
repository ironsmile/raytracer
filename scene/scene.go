package scene

import (
	"fmt"
	"math"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/primitive"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
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

func (s *Scene) Intersect(ray geometry.Ray) (prim primitive.Primitive, retdist float64, normal geometry.Vector) {
	retdist = math.MaxFloat64

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

func (s *Scene) IntersectBBoxEdge(ray geometry.Ray, maxDist float64) bool {
	for _, pr := range s.Primitives {
		if pr.IntersectBBoxEdge(ray, maxDist) {
			return true
		}
	}
	return false
}

func (s *Scene) InitScene() {
	s.Primitives = make([]primitive.Primitive, 0)
	s.Lights = make([]primitive.Primitive, 0)

	rect := primitive.NewRectangle(0.5, 1)
	rect.Name = "rect-floor"
	rect.Mat.Refl = 0
	rect.Mat.Diff = 0.95
	rect.Mat.Color = geometry.NewColor(0.4, 0.3, 0.3)
	rect.SetTransform(
		transform.RotateX(90).Multiply(
			transform.Scale(68, 68, 1),
		).Multiply(
			transform.Translate(geometry.NewVector(0, 0, 5)),
		),
	)

	s.Primitives = append(s.Primitives, rect)

	rect = primitive.NewRectangle(0.5, 1)
	rect.Name = "rect-ceiling"
	rect.Mat.Refl = 0
	rect.Mat.Diff = 0.95
	rect.Mat.Color = geometry.NewColor(0.4, 0.3, 0.3)
	rect.SetTransform(
		transform.RotateX(270).Multiply(
			transform.Scale(68, 68, 1),
		).Multiply(
			transform.Translate(geometry.NewVector(0, 0, 16)),
		),
	)

	s.Primitives = append(s.Primitives, rect)

	rect = primitive.NewRectangle(1, 0.5)
	rect.Name = "rect-left"
	rect.Mat.Refl = 0
	rect.Mat.Diff = 0.95
	rect.Mat.Color = geometry.NewColor(0.4, 0.3, 0.3)
	rect.SetTransform(
		transform.RotateY(270).Multiply(
			transform.Scale(68, 68, 1).Multiply(
				transform.Translate(geometry.NewVector(0, 0, 25)),
			),
		),
	)

	s.Primitives = append(s.Primitives, rect)

	rect = primitive.NewRectangle(1, 0.5)
	rect.Name = "rect-right-mirror"
	rect.Mat.Refl = 1.0
	rect.Mat.Diff = 0.4
	rect.Mat.Color = geometry.NewColor(0.4, 0.3, 0.3)
	rect.SetTransform(
		transform.RotateY(90).Multiply(
			transform.Scale(68, 68, 1).Multiply(
				transform.Translate(geometry.NewVector(0, 0, 10)),
			),
		),
	)

	s.Primitives = append(s.Primitives, rect)

	rect = primitive.NewRectangle(1, 0.5)
	rect.Name = "rect-front"
	rect.Mat.Refl = 0
	rect.Mat.Diff = 0.95
	rect.Mat.Color = geometry.NewColor(0.4, 0.3, 0.3)
	rect.SetTransform(
		transform.Scale(68, 68, 1).Multiply(
			transform.Translate(geometry.NewVector(0, 0, 30)),
		),
	)

	s.Primitives = append(s.Primitives, rect)

	rect = primitive.NewRectangle(1, 0.5)
	rect.Name = "rect-back"
	rect.Mat.Refl = 0
	rect.Mat.Diff = 0.95
	rect.Mat.Color = geometry.NewColor(0.4, 0.3, 0.3)
	rect.SetTransform(
		transform.RotateY(180).Multiply(
			transform.Scale(68, 68, 1).Multiply(
				transform.Translate(geometry.NewVector(0, 0, 25)),
			),
		),
	)

	s.Primitives = append(s.Primitives, rect)

	sphere := primitive.NewSphere(2.5)
	sphere.Name = "big sphere"
	sphere.Mat.Refl = 0.0
	sphere.Mat.Diff = 0.9
	sphere.Mat.Color = geometry.NewColor(1, 0, 0)
	sphere.SetTransform(transform.Translate(geometry.NewVector(1, -0.8, 3)))

	s.Primitives = append(s.Primitives, sphere)

	sphere = primitive.NewSphere(2)
	sphere.Name = "small sphere"
	sphere.Mat.Refl = 0.0
	sphere.Mat.Diff = 0.4
	sphere.Mat.Color = geometry.NewColor(0.7, 0.7, 1)
	sphere.SetTransform(transform.Translate(geometry.NewVector(-5.5, -0.5, 7)))

	s.Primitives = append(s.Primitives, sphere)

	sphere = primitive.NewSphere(1.5)
	sphere.Name = "small sphere far away"
	sphere.Mat.Refl = 0.9
	sphere.Mat.Diff = 0.4
	sphere.Mat.Color = geometry.NewColor(0.5, 1, 0)
	sphere.SetTransform(transform.Translate(geometry.NewVector(-6.5, -2.5, 25)))

	s.Primitives = append(s.Primitives, sphere)

	triangle := primitive.NewTriangle([3]geometry.Point{
		*geometry.NewPoint(-10.99, 3, 0),  // a
		*geometry.NewPoint(-10.99, 0, -3), // b
		*geometry.NewPoint(-10.99, 0, 3),  // c
	})
	triangle.Name = "Green triangle"
	triangle.Mat.Refl = 0.0
	triangle.Mat.Diff = 0.3
	triangle.Mat.Color = geometry.NewColor(0.3, 1, 0)

	s.Primitives = append(s.Primitives, triangle)

	sphere = primitive.NewSphere(0.1)
	sphere.Name = "Visible light source"
	sphere.Light = true
	sphere.LightSource = *geometry.NewPoint(0, 5, 5)
	sphere.Mat.Color = geometry.NewColor(0.9, 0.9, 0.9)
	sphere.SetTransform(transform.Translate(sphere.LightSource.Vector()))

	s.Primitives = append(s.Primitives, sphere)
	s.Lights = append(s.Lights, sphere)

	sphere = primitive.NewSphere(0.1)
	sphere.Name = "Invisible lightsource"
	sphere.Light = true
	sphere.LightSource = *geometry.NewPoint(2, 5, 1)
	sphere.Mat.Color = geometry.NewColor(0.9, 0.9, 0.9)
	sphere.SetTransform(transform.Translate(sphere.LightSource.Vector()))

	s.Primitives = append(s.Primitives, sphere)
	s.Lights = append(s.Lights, sphere)

	sphere = primitive.NewSphere(0.1)
	sphere.Name = "Behid the shoulder lightsource"
	sphere.Light = true
	sphere.LightSource = *geometry.NewPoint(2, 5, -10)
	sphere.Mat.Color = geometry.NewColor(0.9, 0.9, 0.9)
	sphere.SetTransform(transform.Translate(sphere.LightSource.Vector()))

	s.Primitives = append(s.Primitives, sphere)
	s.Lights = append(s.Lights, sphere)

	teaPotCenter := geometry.NewVector(-3, 0, 5)
	if teapot, err := primitive.NewObject("data/objs/teapot.obj"); err != nil {
		fmt.Printf("Error loading obj teapot: %s\n", err)
	} else {
		teapot.Name = "First teapod"
		teapot.Mat.Refl = 0.0
		teapot.Mat.Diff = 0.3
		teapot.Mat.Color = geometry.NewColor(0.3, 1, 0)
		teapot.SetTransform(
			transform.Translate(teaPotCenter).Multiply(transform.UniformScale(0.01)),
		)

		s.Primitives = append(s.Primitives, teapot)
	}

	blueRect := primitive.NewRectangle(1, 0.5)
	blueRect.Name = "Blue Rectangle"
	blueRect.Mat.Color = geometry.NewColor(0, 0, 1)
	blueRect.Mat.Refl = 0.5
	blueRect.Mat.Diff = 0.8
	blueRect.SetTransform(
		transform.Translate(geometry.NewVector(-10, 0, 0)).Multiply(transform.RotateY(-90)),
	)
	s.Primitives = append(s.Primitives, blueRect)
}

func NewScene() *Scene {
	scn := new(Scene)
	return scn
}
