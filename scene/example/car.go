package example

import (
	"fmt"
	"path/filepath"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"
	"github.com/ironsmile/raytracer/primitive"
	"github.com/ironsmile/raytracer/transform"
)

// GetCarScene returns a predominantly emtpy scene with the alfa147 in the middle
func GetCarScene() ([]primitive.Primitive, []primitive.Primitive) {
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

	// "Visible light source"
	sphere := primitive.NewSphere(0.1)
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(0, 5, 5)
	sphere.Shape().SetMaterial(mat.Material{
		Color: geometry.NewColor(0.9, 0.9, 0.9),
	})
	sphere.SetTransform(transform.Translate(sphere.LightSource))

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

	primitives = append(primitives, sphere)
	lights = append(lights, sphere)

	// "Behid the shoulder lightsource"
	sphere = primitive.NewSphere(0.1)
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(2, 5, -10)
	sphere.Shape().SetMaterial(mat.Material{
		Color: geometry.NewColor(0.9, 0.9, 0.9),
	})
	sphere.SetTransform(transform.Translate(sphere.LightSource))

	primitives = append(primitives, sphere)
	lights = append(lights, sphere)

	alfaPath := filepath.Join("data", "objs", "alfa147.obj")
	if obj, err := primitive.NewObject(alfaPath); err != nil {
		fmt.Printf("Error loading obj alfa147: %s\n", err)
	} else {
		objTransform := transform.Translate(geometry.NewVector(-2.5, -5, 3)).Multiply(
			transform.UniformScale(1).Multiply(
				transform.RotateX(-180),
			),
		)
		obj.SetTransform(objTransform)

		primitives = append(primitives, obj)
	}

	return primitives, lights
}
