package example

import (
	"fmt"
	"path/filepath"

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
	rect := primitive.NewQuad(
		geometry.NewVector(-25, -5, 30),
		geometry.NewVector(10, -5, 30),
		geometry.NewVector(10, -5, -25),
		geometry.NewVector(-25, -5, -25),
	)
	rect.Shape().SetMaterial(wallMaterial)
	primitive.SetName(rect.GetID(), "rect-floor")
	primitives = append(primitives, rect)

	// "rect-ceiling"
	rect = primitive.NewQuad(
		geometry.NewVector(-25, 16, 30),
		geometry.NewVector(-25, 16, -25),
		geometry.NewVector(10, 16, -25),
		geometry.NewVector(10, 16, 30),
	)
	rect.Shape().SetMaterial(wallMaterial)
	primitive.SetName(rect.GetID(), "rect-ceiling")
	primitives = append(primitives, rect)

	// "rect-left"
	rect = primitive.NewQuad(
		geometry.NewVector(-25, 16, -25),
		geometry.NewVector(-25, 16, 30),
		geometry.NewVector(-25, -5, 30),
		geometry.NewVector(-25, -5, -25),
	)
	rect.Shape().SetMaterial(wallMaterial)
	primitive.SetName(rect.GetID(), "rect-left")
	primitives = append(primitives, rect)

	// "rect-right-mirror"
	rect = primitive.NewQuad(
		geometry.NewVector(10, 16, -25),
		geometry.NewVector(10, -5, -25),
		geometry.NewVector(10, -5, 30),
		geometry.NewVector(10, 16, 30),
	)
	rect.Shape().SetMaterial(reflectiveWallMaterial)
	primitive.SetName(rect.GetID(), "rect-right-mirror")
	primitives = append(primitives, rect)

	// "rect-front"
	rect = primitive.NewQuad(
		geometry.NewVector(-25, 16, 30),
		geometry.NewVector(10, 16, 30),
		geometry.NewVector(10, -5, 30),
		geometry.NewVector(-25, -5, 30),
	)
	rect.Shape().SetMaterial(wallMaterial)
	primitive.SetName(rect.GetID(), "rect-front")
	primitives = append(primitives, rect)

	// "rect-back"
	rect = primitive.NewQuad(
		geometry.NewVector(-25, 16, -25),
		geometry.NewVector(-25, -5, -25),
		geometry.NewVector(10, -5, -25),
		geometry.NewVector(10, 16, -25),
	)
	rect.Shape().SetMaterial(wallMaterial)
	primitive.SetName(rect.GetID(), "rect-back")
	primitives = append(primitives, rect)

	// "big sphere"
	sphere := primitive.NewSphere(2.5)
	sphere.Shape().SetMaterial(mat.Material{
		Refl:  0.0,
		Diff:  0.9,
		Color: geometry.NewColor(1, 0, 0),
	})
	sphere.SetTransform(transform.Translate(geometry.NewVector(1, -0.8, 3)))
	primitive.SetName(sphere.GetID(), "big red sphere")
	primitives = append(primitives, sphere)

	// "small sphere"
	sphere = primitive.NewSphere(2)
	sphere.Shape().SetMaterial(mat.Material{
		Refl:      0.0,
		Refr:      0.6,
		RefrIndex: 1.5,
		Diff:      0.4,
		Color:     geometry.NewColor(0.7, 0.7, 1),
	})
	sphere.SetTransform(transform.Translate(geometry.NewVector(-5.5, -0.5, 7)))
	primitive.SetName(sphere.GetID(), "small sphere")
	primitives = append(primitives, sphere)

	// "small sphere far away"
	sphere = primitive.NewSphere(1.5)
	sphere.Shape().SetMaterial(mat.Material{
		Refl:  0.9,
		Diff:  0.4,
		Color: geometry.NewColor(0.5, 1, 0),
	})
	sphere.SetTransform(transform.Translate(geometry.NewVector(-6.5, -2.5, 25)))
	primitive.SetName(sphere.GetID(), "small sphere far away")
	primitives = append(primitives, sphere)

	// "Green triangle"
	triangle := primitive.NewTriangle([3]geometry.Vector{
		geometry.NewVector(-10.99, 3, 0),  // a
		geometry.NewVector(-10.99, 0, -3), // b
		geometry.NewVector(-10.99, 0, 3),  // c
	})
	triangle.Shape().SetMaterial(mat.Material{
		Refl:  0.0,
		Diff:  0.3,
		Color: geometry.NewColor(0.3, 1, 0),
	})
	primitive.SetName(triangle.GetID(), "green triangle")
	primitives = append(primitives, triangle)

	// "Visible light source"
	sphere = primitive.NewSphere(0.1)
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(0, 5, 5)
	sphere.Shape().SetMaterial(mat.Material{
		Color: geometry.NewColor(0.9, 0.9, 0.9),
	})
	sphere.SetTransform(transform.Translate(sphere.LightSource))
	primitive.SetName(sphere.GetID(), "Visible light source")

	primitives = append(primitives, sphere)
	lights = append(lights, sphere)

	// "Invisible lightsource"
	sphere = primitive.NewSphere(0.1)
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(2, 5, 1)
	sphere.Shape().SetMaterial(mat.Material{
		Color: geometry.NewColor(0.9, 0.9, 0.9),
	})
	sphere.SetTransform(transform.Translate(sphere.LightSource))
	primitive.SetName(sphere.GetID(), "Invisible light source")

	primitives = append(primitives, sphere)
	lights = append(lights, sphere)

	// "Behind the shoulder lightsource"
	sphere = primitive.NewSphere(0.1)
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(2, 5, -10)
	sphere.Shape().SetMaterial(mat.Material{
		Color: geometry.NewColor(0.9, 0.9, 0.9),
	})
	sphere.SetTransform(transform.Translate(sphere.LightSource))
	primitive.SetName(sphere.GetID(), "Behind the shoulder lightsource")

	primitives = append(primitives, sphere)
	lights = append(lights, sphere)

	teapotPath := filepath.Join("data", "objs", "teapot.obj")
	if obj, err := primitive.NewObject(teapotPath); err != nil {
		fmt.Printf("error loading obj teapot: %s\n", err)
	} else {
		objTransform := transform.Translate(geometry.NewVector(-3, 0, 5)).Multiply(
			transform.UniformScale(0.2),
		)
		obj.SetTransform(objTransform)
		primitive.SetName(obj.GetID(), "teapot")

		primitives = append(primitives, obj)
	}

	// "Blue Rectangle"
	quad := primitive.NewQuad(
		geometry.NewVector(-1, 0.5, 0),
		geometry.NewVector(1, 0.5, 0),
		geometry.NewVector(1, -0.5, 0),
		geometry.NewVector(-1, -0.5, 0),
	)
	quad.Shape().SetMaterial(mat.Material{
		Refl:  0.5,
		Diff:  0.8,
		Color: geometry.NewColor(0, 0, 1),
	})
	quad.SetTransform(
		transform.Translate(geometry.NewVector(-10, 0, 0)).Multiply(
			transform.RotateY(-90),
		),
	)
	primitive.SetName(quad.GetID(), "Blue Rectangle")
	primitives = append(primitives, quad)

	// Cyan cylinder
	cyl := primitive.NewCylinder(0.5, geometry.NewVector(0, 0, 0), geometry.NewVector(0, 2, 0))
	cyl.Shape().SetMaterial(mat.Material{
		Diff:  1,
		Color: geometry.NewColor(0, 1, 1),
	})
	cyl.SetTransform(
		transform.Translate(geometry.NewVector(
			-4, -5, -1,
		)),
	)
	primitive.SetName(cyl.GetID(), "Cyan Cylinder")
	primitives = append(primitives, cyl)

	return primitives, lights
}
