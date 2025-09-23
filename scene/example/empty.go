package example

import (
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"
	"github.com/ironsmile/raytracer/primitive"
	"github.com/ironsmile/raytracer/transform"
)

// GetEmptyScene returs a scene with no primitives and one light.
func GetEmptyScene() ([]primitive.Primitive, []primitive.Primitive) {
	var primitives []primitive.Primitive
	var lights []primitive.Primitive

	// "Behid the shoulder lightsource"
	sphere := primitive.NewSphere(0.1)
	sphere.Light = true
	sphere.LightSource = geometry.NewVector(2, 5, -10)
	sphere.Shape().SetMaterial(mat.Material{
		Color: geometry.NewColor(0.9, 0.9, 0.9),
	})
	sphere.SetTransform(transform.Translate(sphere.LightSource))

	primitives = append(primitives, sphere)
	lights = append(lights, sphere)

	return primitives, lights
}
