package camera

import (
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/transform"
)

type ProjectiveCamera struct {
	CameraToWorld             *transform.Transform
	ShutterOpen, ShutterClose float64

	cameraToScreen *transform.Transform
	rasterToCamera *transform.Transform
	screenToRaster *transform.Transform
	rasterToScreen *transform.Transform

	lensRadius    float64
	focalDistance float64
}

func NewProjectiveCamera(
	camToWorld *transform.Transform,
	proj *transform.Transform,
	screenWindow [4]float64,
	sopen, sclose, lensr, focald float64,
	width, height float64,
) *ProjectiveCamera {
	cam := &ProjectiveCamera{
		CameraToWorld:  camToWorld,
		ShutterOpen:    sopen,
		ShutterClose:   sclose,
		cameraToScreen: proj,
	}

	cam.lensRadius = lensr
	cam.focalDistance = focald

	cam.screenToRaster = transform.Scale(width, height, 1.0).
		Multiply(transform.Scale(1.0/(screenWindow[1]-screenWindow[0]),
			1.0/(screenWindow[2]-screenWindow[3]), 1.0)).
		Multiply(transform.Translate(geometry.NewVector(-screenWindow[0],
			-screenWindow[3], 0.0)))

	cam.rasterToScreen = cam.screenToRaster.Inverse()

	cam.rasterToCamera = cam.cameraToScreen.Inverse().Multiply(cam.rasterToScreen)

	return cam
}
