package example

import (
	"fmt"

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
		geometry.NewVector(-0.5, 1, 0),
		geometry.NewVector(0.5, 1, 0),
		geometry.NewVector(0.5, -1, 0),
		geometry.NewVector(-0.5, -1, 0),
	)
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
	rect = primitive.NewQuad(
		geometry.NewVector(-0.5, 1, 0),
		geometry.NewVector(0.5, 1, 0),
		geometry.NewVector(0.5, -1, 0),
		geometry.NewVector(-0.5, -1, 0),
	)
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
	rect = primitive.NewQuad(
		geometry.NewVector(-1, 0.5, 0),
		geometry.NewVector(1, 0.5, 0),
		geometry.NewVector(1, -0.5, 0),
		geometry.NewVector(-1, -0.5, 0),
	)
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
	rect = primitive.NewQuad(
		geometry.NewVector(-1, 0.5, 0),
		geometry.NewVector(1, 0.5, 0),
		geometry.NewVector(1, -0.5, 0),
		geometry.NewVector(-1, -0.5, 0),
	)
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
	rect = primitive.NewQuad(
		geometry.NewVector(-1, 0.5, 0),
		geometry.NewVector(1, 0.5, 0),
		geometry.NewVector(1, -0.5, 0),
		geometry.NewVector(-1, -0.5, 0),
	)
	rect.Mat = &wallMaterial
	rect.SetTransform(
		transform.Scale(68, 68, 1).Multiply(
			transform.Translate(geometry.NewVector(0, 0, 30)),
		),
	)
	primitive.SetName(rect.GetID(), "rect-front")

	primitives = append(primitives, rect)

	// "rect-back"
	rect = primitive.NewQuad(
		geometry.NewVector(-1, 0.5, 0),
		geometry.NewVector(1, 0.5, 0),
		geometry.NewVector(1, -0.5, 0),
		geometry.NewVector(-1, -0.5, 0),
	)
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

	// "Visible light source"
	sphere := primitive.NewSphere(0.1)
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(0, 5, 5)
	sphere.Mat = &mat.Material{
		Color: geometry.NewColor(0.9, 0.9, 0.9),
	}
	sphere.SetTransform(transform.Translate(sphere.LightSource))

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

	primitives = append(primitives, sphere)
	lights = append(lights, sphere)

	// "Behid the shoulder lightsource"
	sphere = primitive.NewSphere(0.1)
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(2, 5, -10)
	sphere.Mat = &mat.Material{
		Color: geometry.NewColor(0.9, 0.9, 0.9),
	}
	sphere.SetTransform(transform.Translate(sphere.LightSource))

	primitives = append(primitives, sphere)
	lights = append(lights, sphere)

	if obj, err := primitive.NewObject("data/objs/alfa147.obj"); err != nil {
		fmt.Printf("Error loading obj alfa147: %s\n", err)
	} else {
		objTransform := transform.Translate(geometry.NewVector(-2.5, -5, 3)).Multiply(
			transform.UniformScale(0.05).Multiply(
				transform.RotateX(-90),
				// Multiply(
				//  transform.RotateZ(140),
				// ),
			),
		)
		obj.Mat = &mat.Material{
			Refl:  0.0,
			Diff:  0.3,
			Color: geometry.NewColor(0.729, 0.572, 0.780),
		}
		obj.SetTransform(objTransform)

		primitives = append(primitives, obj)
	}

	return primitives, lights
}
