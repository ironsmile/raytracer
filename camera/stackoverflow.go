package camera

import (
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/geometry"
)

type StackOverflowCamera struct {
	Film film.Film

	position  *geometry.Point
	direction *geometry.Vector
	up        *geometry.Vector
	right     *geometry.Vector

	rasterW, rasterH float64
}

func (soc *StackOverflowCamera) GenerateRay(screenX, screenY float64) (*geometry.Ray,
	float64) {

	normalizedI := (screenX / soc.rasterW) - 0.5
	normalizedJ := (screenY / soc.rasterH) - 0.5

	imagePoint := soc.position.PlusVector(soc.right.MultiplyScalar(normalizedI).Plus(
		soc.up.MultiplyScalar(normalizedJ)).Plus(soc.direction))

	rayDirection := imagePoint.Minus(soc.position)

	return geometry.NewRay(*soc.position, *rayDirection), 1.0
}

func NewStackOverflowCamera(camPosition, camLookAtPoint *geometry.Point,
	camUp *geometry.Vector, f film.Film) *StackOverflowCamera {
	cam := &StackOverflowCamera{Film: f, position: camPosition}

	cam.rasterW = float64(f.Width())
	cam.rasterH = float64(f.Height())

	cam.direction = geometry.Normalize(camLookAtPoint.Minus(camPosition))
	cam.right = cam.direction.Cross(camUp)
	cam.up = cam.right.Cross(cam.direction)

	return cam
}
