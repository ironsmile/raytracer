package scene

import (
	"fmt"

	"github.com/ironsmile/raytracer/mat"

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

	wallMaterial := mat.Material{}
	wallMaterial.Refl = 0
	wallMaterial.Diff = 0.95
	wallMaterial.Color = geometry.NewColor(0.4, 0.3, 0.3)

	reflectiveWallMaterial := mat.Material{}
	reflectiveWallMaterial.Refl = 1.0
	reflectiveWallMaterial.Diff = 0.4
	reflectiveWallMaterial.Color = geometry.NewColor(0.4, 0.3, 0.3)

	// "rect-floor"
	rect := primitive.NewRectangle(0.5, 1)
	rect.Mat = &wallMaterial
	rect.SetTransform(
		transform.Translate(geometry.NewVector(-10, -5, 0)).Multiply(
			transform.RotateX(90).Multiply(
				transform.Scale(90, 68, 1),
			),
		),
	)

	s.Primitives = append(s.Primitives, rect)

	// "rect-ceiling"
	rect = primitive.NewRectangle(0.5, 1)
	rect.Mat = &wallMaterial
	rect.SetTransform(
		transform.Translate(geometry.NewVector(-10, 16, 0)).Multiply(
			transform.RotateX(270).Multiply(
				transform.Scale(90, 68, 1),
			),
		),
	)

	s.Primitives = append(s.Primitives, rect)

	// "rect-left"
	rect = primitive.NewRectangle(1, 0.5)
	rect.Mat = &wallMaterial
	rect.SetTransform(
		transform.RotateY(270).Multiply(
			transform.Scale(68, 68, 1).Multiply(
				transform.Translate(geometry.NewVector(0, 0, 25)),
			),
		),
	)

	s.Primitives = append(s.Primitives, rect)

	// "rect-right-mirror"
	rect = primitive.NewRectangle(1, 0.5)
	rect.Mat = &reflectiveWallMaterial
	rect.SetTransform(
		transform.RotateY(90).Multiply(
			transform.Scale(68, 68, 1).Multiply(
				transform.Translate(geometry.NewVector(0, 0, 10)),
			),
		),
	)

	s.Primitives = append(s.Primitives, rect)

	// "rect-front"
	rect = primitive.NewRectangle(1, 0.5)
	rect.Mat = &wallMaterial
	rect.SetTransform(
		transform.Scale(68, 68, 1).Multiply(
			transform.Translate(geometry.NewVector(0, 0, 30)),
		),
	)

	s.Primitives = append(s.Primitives, rect)

	// "rect-back"
	rect = primitive.NewRectangle(1, 0.5)
	rect.Mat = &wallMaterial
	rect.SetTransform(
		transform.RotateY(180).Multiply(
			transform.Scale(68, 68, 1).Multiply(
				transform.Translate(geometry.NewVector(0, 0, 25)),
			),
		),
	)

	s.Primitives = append(s.Primitives, rect)

	// "big sphere"
	sphere := primitive.NewSphere(2.5)
	sphere.Mat = &mat.Material{
		Refl:  0.0,
		Diff:  0.9,
		Color: geometry.NewColor(1, 0, 0),
	}
	sphere.SetTransform(transform.Translate(geometry.NewVector(1, -0.8, 3)))

	s.Primitives = append(s.Primitives, sphere)

	// "small sphere"
	sphere = primitive.NewSphere(2)
	sphere.Mat = &mat.Material{
		Refl:  0.0,
		Diff:  0.4,
		Color: geometry.NewColor(0.7, 0.7, 1),
	}
	sphere.SetTransform(transform.Translate(geometry.NewVector(-5.5, -0.5, 7)))

	s.Primitives = append(s.Primitives, sphere)

	// "small sphere far away"
	sphere = primitive.NewSphere(1.5)
	sphere.Mat = &mat.Material{
		Refl:  0.9,
		Diff:  0.4,
		Color: geometry.NewColor(0.5, 1, 0),
	}
	sphere.SetTransform(transform.Translate(geometry.NewVector(-6.5, -2.5, 25)))

	s.Primitives = append(s.Primitives, sphere)

	triangle := primitive.NewTriangle([3]geometry.Vector{
		geometry.NewVector(-10.99, 3, 0),  // a
		geometry.NewVector(-10.99, 0, -3), // b
		geometry.NewVector(-10.99, 0, 3),  // c
		// "Green triangle"
	})
	triangle.Mat = &mat.Material{
		Refl:  0.0,
		Diff:  0.3,
		Color: geometry.NewColor(0.3, 1, 0),
	}

	s.Primitives = append(s.Primitives, triangle)

	// "Visible light source"
	sphere = primitive.NewSphere(0.1)
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(0, 5, 5)
	sphere.Mat = &mat.Material{
		Color: geometry.NewColor(0.9, 0.9, 0.9),
	}
	sphere.SetTransform(transform.Translate(sphere.LightSource))

	s.Primitives = append(s.Primitives, sphere)
	s.Lights = append(s.Lights, sphere)

	// "Invisible lightsource"
	sphere = primitive.NewSphere(0.1)
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(2, 5, 1)
	sphere.Mat = &mat.Material{
		Color: geometry.NewColor(0.9, 0.9, 0.9),
	}
	sphere.SetTransform(transform.Translate(sphere.LightSource))

	s.Primitives = append(s.Primitives, sphere)
	s.Lights = append(s.Lights, sphere)

	// "Behid the shoulder lightsource"
	sphere = primitive.NewSphere(0.1)
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(2, 5, -10)
	sphere.Mat = &mat.Material{
		Color: geometry.NewColor(0.9, 0.9, 0.9),
	}
	sphere.SetTransform(transform.Translate(sphere.LightSource))

	s.Primitives = append(s.Primitives, sphere)
	s.Lights = append(s.Lights, sphere)

	// if obj, err := primitive.NewObject("data/objs/alfa147.obj"); err != nil {
	// 	fmt.Printf("Error loading obj alfa147: %s\n", err)
	// } else {
	// 	objTransform := transform.Translate(geometry.NewVector(0, -3, 1)).Multiply(
	// 		transform.UniformScale(0.05).Multiply(
	// 			transform.RotateX(-90),
	// 			// Multiply(
	// 			// 	transform.RotateZ(140),
	// 			// ),
	// 		),
	// 	)
	// 	obj.Mat = &mat.Material{
	// 		Refl:  0.0,
	// 		Diff:  0.3,
	// 		Color: geometry.NewColor(0.557, 0.286, 0.643),
	// 	}
	// 	obj.SetTransform(objTransform)

	// 	s.Primitives = append(s.Primitives, obj)
	// }

	if obj, err := primitive.NewObject("data/objs/teapot.obj"); err != nil {
		fmt.Printf("Error loading obj teapot: %s\n", err)
	} else {
		objTransform := transform.Translate(geometry.NewVector(-3, 0, 5)).Multiply(
			transform.UniformScale(0.01),
		)
		obj.Mat = &mat.Material{
			Refl:  0.0,
			Diff:  0.3,
			Color: geometry.NewColor(0.557, 0.286, 0.643),
		}
		obj.SetTransform(objTransform)

		s.Primitives = append(s.Primitives, obj)
	}

	// "Blue Rectangle"
	blueRect := primitive.NewRectangle(1, 0.5)
	blueRect.Mat = &mat.Material{
		Refl:  0.5,
		Diff:  0.8,
		Color: geometry.NewColor(0, 0, 1),
	}
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
