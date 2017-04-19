package scene

import (
	"fmt"

	"github.com/ironsmile/raytracer/accel"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/primitive"
	"github.com/ironsmile/raytracer/transform"
)

// Scene is a type which is responsible for loading and managing a scene for rendering.
type Scene struct {
	Primitives []primitive.Primitive
	Lights     []primitive.Primitive

	accel primitive.Primitive
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
func (s *Scene) InitScene() {
	s.Primitives = make([]primitive.Primitive, 0)
	s.Lights = make([]primitive.Primitive, 0)

	rect := primitive.NewRectangle(0.5, 1)
	rect.Name = "rect-floor"
	rect.Mat.Refl = 0
	rect.Mat.Diff = 0.95
	rect.Mat.Color = geometry.NewColor(0.4, 0.3, 0.3)
	rect.SetTransform(
		transform.Translate(geometry.NewVector(-10, -5, 0)).Multiply(
			transform.RotateX(90).Multiply(
				transform.Scale(90, 68, 1),
			),
		),
	)

	s.Primitives = append(s.Primitives, rect)

	rect = primitive.NewRectangle(0.5, 1)
	rect.Name = "rect-ceiling"
	rect.Mat.Refl = 0
	rect.Mat.Diff = 0.95
	rect.Mat.Color = geometry.NewColor(0.4, 0.3, 0.3)
	rect.SetTransform(
		transform.Translate(geometry.NewVector(-10, 16, 0)).Multiply(
			transform.RotateX(270).Multiply(
				transform.Scale(90, 68, 1),
			),
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

	triangle := primitive.NewTriangle([3]geometry.Vector{
		geometry.NewVector(-10.99, 3, 0),  // a
		geometry.NewVector(-10.99, 0, -3), // b
		geometry.NewVector(-10.99, 0, 3),  // c
	})
	triangle.Name = "Green triangle"
	triangle.Mat.Refl = 0.0
	triangle.Mat.Diff = 0.3
	triangle.Mat.Color = geometry.NewColor(0.3, 1, 0)

	s.Primitives = append(s.Primitives, triangle)

	sphere = primitive.NewSphere(0.1)
	sphere.Name = "Visible light source"
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(0, 5, 5)
	sphere.Mat.Color = geometry.NewColor(0.9, 0.9, 0.9)
	sphere.SetTransform(transform.Translate(sphere.LightSource))

	s.Primitives = append(s.Primitives, sphere)
	s.Lights = append(s.Lights, sphere)

	sphere = primitive.NewSphere(0.1)
	sphere.Name = "Invisible lightsource"
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(2, 5, 1)
	sphere.Mat.Color = geometry.NewColor(0.9, 0.9, 0.9)
	sphere.SetTransform(transform.Translate(sphere.LightSource))

	s.Primitives = append(s.Primitives, sphere)
	s.Lights = append(s.Lights, sphere)

	sphere = primitive.NewSphere(0.1)
	sphere.Name = "Behid the shoulder lightsource"
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(2, 5, -10)
	sphere.Mat.Color = geometry.NewColor(0.9, 0.9, 0.9)
	sphere.SetTransform(transform.Translate(sphere.LightSource))

	s.Primitives = append(s.Primitives, sphere)
	s.Lights = append(s.Lights, sphere)

	// if obj, err := primitive.NewObject("data/objs/alfa147.obj"); err != nil {
	// 	fmt.Printf("Error loading obj alfa147: %s\n", err)
	// } else {
	// 	obj.Name = "First alfa147"
	// 	obj.Mat.Refl = 0.0
	// 	obj.Mat.Diff = 0.3
	// 	obj.Mat.Color = geometry.NewColor(0.3, 1, 0)
	// 	obj.SetTransform(
	// 		transform.Translate(geometry.NewVector(0, -3, 1)).Multiply(
	// 			transform.UniformScale(0.05).Multiply(
	// 				transform.RotateX(-90),
	// 				// Multiply(
	// 				// 	transform.RotateZ(140),
	// 				// ),
	// 			),
	// 		),
	// 	)

	// 	var prims []primitive.Primitive

	// 	for _, objShape := range obj.Shape().GetAllShapes() {
	// 		prims = append(prims, primitive.FromShape(objShape))
	// 	}

	// 	// s.Primitives = append(s.Primitives, accel.NewGrid(prims))
	// 	s.Primitives = append(s.Primitives, obj)
	// }

	if obj, err := primitive.NewObject("data/objs/teapot.obj"); err != nil {
		fmt.Printf("Error loading obj teapot: %s\n", err)
	} else {
		obj.Name = "First teapot"
		obj.Mat.Refl = 0.0
		obj.Mat.Diff = 0.3
		obj.Mat.Color = geometry.NewColor(0.3, 1, 0)
		obj.SetTransform(
			transform.Translate(geometry.NewVector(-3, 0, 5)).Multiply(
				transform.UniformScale(0.01),
			),
		)

		var prims []primitive.Primitive

		for _, objShape := range obj.Shape().GetAllShapes() {
			prims = append(prims, primitive.FromShape(objShape))
		}

		// s.Primitives = append(s.Primitives, accel.NewGrid(prims))
		s.Primitives = append(s.Primitives, obj)
	}

	blueRect := primitive.NewRectangle(1, 0.5)
	blueRect.Name = "Blue Rectangle"
	blueRect.Mat.Color = geometry.NewColor(0, 0, 1)
	blueRect.Mat.Refl = 0.5
	blueRect.Mat.Diff = 0.8
	blueRect.SetTransform(
		transform.Translate(geometry.NewVector(-10, 0, 0)).Multiply(
			transform.RotateY(-90),
		),
	)
	s.Primitives = append(s.Primitives, blueRect)

	s.accel = accel.NewGrid(s.Primitives)
}

// NewScene returns a new demo scene
func NewScene() *Scene {
	scn := new(Scene)
	return scn
}
