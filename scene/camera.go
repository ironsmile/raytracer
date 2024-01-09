package scene

import (
	"github.com/ironsmile/raytracer/camera"
	"github.com/ironsmile/raytracer/geometry"
)

func GetCamera(w, h float64) camera.Camera {
	pos := geometry.NewVector(0, 0, -5)
	lookAtPoint := geometry.NewVector(0, 0, 1)
	up := geometry.NewVector(0, 1, 0)

	return camera.NewPinhole(pos, lookAtPoint, up, 1, w, h)
}
