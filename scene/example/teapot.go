package example

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"
	"github.com/ironsmile/raytracer/primitive"
	"github.com/ironsmile/raytracer/transform"
)

// GetTeapotScene returns the default teapot scene used throughout the development
func GetTeapotScene() ([]primitive.Primitive, []primitive.Primitive) {
	var primitives []primitive.Primitive
	var lights []primitive.Primitive

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
	primitive.SetName(rect.GetID(), "rect-floor")

	primitives = append(primitives, rect)

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
	primitive.SetName(rect.GetID(), "rect-ceiling")

	primitives = append(primitives, rect)

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
	primitive.SetName(rect.GetID(), "rect-left")

	primitives = append(primitives, rect)

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
	primitive.SetName(rect.GetID(), "rect-right-mirror")

	primitives = append(primitives, rect)

	// "rect-front"
	rect = primitive.NewRectangle(1, 0.5)
	rect.Mat = &wallMaterial
	rect.SetTransform(
		transform.Scale(68, 68, 1).Multiply(
			transform.Translate(geometry.NewVector(0, 0, 30)),
		),
	)
	primitive.SetName(rect.GetID(), "rect-front")

	primitives = append(primitives, rect)

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
	primitive.SetName(rect.GetID(), "rect-back")

	primitives = append(primitives, rect)

	// "big sphere"
	sphere := primitive.NewSphere(2.5)
	sphere.Mat = &mat.Material{
		Refl:  0.0,
		Diff:  0.9,
		Color: geometry.NewColor(1, 0, 0),
	}
	sphere.SetTransform(transform.Translate(geometry.NewVector(1, -0.8, 3)))
	primitive.SetName(sphere.GetID(), "big red sphere")

	primitives = append(primitives, sphere)

	// "small sphere"
	sphere = primitive.NewSphere(2)
	sphere.Mat = &mat.Material{
		Refl:  0.0,
		Diff:  0.4,
		Color: geometry.NewColor(0.7, 0.7, 1),
	}
	sphere.SetTransform(transform.Translate(geometry.NewVector(-5.5, -0.5, 7)))
	primitive.SetName(sphere.GetID(), "small sphere")

	primitives = append(primitives, sphere)

	// "small sphere far away"
	sphere = primitive.NewSphere(1.5)
	sphere.Mat = &mat.Material{
		Refl:  0.9,
		Diff:  0.4,
		Color: geometry.NewColor(0.5, 1, 0),
	}
	sphere.SetTransform(transform.Translate(geometry.NewVector(-6.5, -2.5, 25)))
	primitive.SetName(sphere.GetID(), "small sphere far away")

	primitives = append(primitives, sphere)

	// "Green triangle"
	triangle := primitive.NewTriangle([3]geometry.Vector{
		geometry.NewVector(-10.99, 3, 0),  // a
		geometry.NewVector(-10.99, 0, -3), // b
		geometry.NewVector(-10.99, 0, 3),  // c
	})
	triangle.Mat = &mat.Material{
		Refl:  0.0,
		Diff:  0.3,
		Color: geometry.NewColor(0.3, 1, 0),
	}
	primitive.SetName(triangle.GetID(), "green triangle")

	primitives = append(primitives, triangle)

	// "Visible light source"
	sphere = primitive.NewSphere(0.1)
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(0, 5, 5)
	sphere.Mat = &mat.Material{
		Color: geometry.NewColor(0.9, 0.9, 0.9),
	}
	sphere.SetTransform(transform.Translate(sphere.LightSource))
	primitive.SetName(sphere.GetID(), "Visible light source")

	primitives = append(primitives, sphere)
	lights = append(lights, sphere)

	// "Invisible lightsource"
	sphere = primitive.NewSphere(0.1)
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(2, 5, 1)
	sphere.Mat = &mat.Material{
		Color: geometry.NewColor(0.9, 0.9, 0.9),
	}
	sphere.SetTransform(transform.Translate(sphere.LightSource))
	primitive.SetName(sphere.GetID(), "Invisible light source")

	primitives = append(primitives, sphere)
	lights = append(lights, sphere)

	// "Behind the shoulder lightsource"
	sphere = primitive.NewSphere(0.1)
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(2, 5, -10)
	sphere.Mat = &mat.Material{
		Color: geometry.NewColor(0.9, 0.9, 0.9),
	}
	sphere.SetTransform(transform.Translate(sphere.LightSource))
	primitive.SetName(sphere.GetID(), "Behind the shoulder lightsource")

	primitives = append(primitives, sphere)
	lights = append(lights, sphere)

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
		primitive.SetName(obj.GetID(), "teapot")

		primitives = append(primitives, obj)
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
	primitive.SetName(blueRect.GetID(), "Blue Rectangle")
	primitives = append(primitives, blueRect)

	return primitives, lights
}
