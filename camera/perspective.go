package camera

import (
	"fmt"
	"math"

	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/transform"
	"github.com/ironsmile/raytracer/utils"
)

type PerspectiveCamera struct {
	ProjectiveCamera

	dxCamera *geometry.Vector
	dyCamera *geometry.Vector
}

func (p *PerspectiveCamera) GenerateRay(screenX, screenY float64) (*geometry.Ray,
	float64) {

	Pras := geometry.NewPoint(screenX, screenY, 0)
	Pcamera := geometry.Normalize(p.rasterToCamera.Point(Pras).Vector())

	if screenX == float64(DEBUG_X) && screenY == float64(DEBUG_Y) {
		fmt.Printf("Camera coords: %s\n", Pcamera)
	}

	ray := geometry.NewRayFull(*geometry.NewPoint(0, 0, 0), *Pcamera, 0.0,
		math.MaxFloat64)

	if p.lensRadius > 0.0 {
		wF := float64(p.Film.Width())
		hF := float64(p.Film.Height())

		sampleLensU := screenX / wF
		sampleLensV := screenY / hF

		lensU, lensV := utils.ConcentricSampleDisk(sampleLensU, sampleLensV)
		lensU *= p.lensRadius
		lensV *= p.lensRadius

		ft := p.focalDistance / ray.Direction.Z
		Pfocus := ray.AtTime(ft)

		ray.Origin = geometry.NewPoint(lensU, lensV, 0.0)
		ray.Direction = geometry.Normalize(Pfocus.Minus(ray.Origin))
	}

	ray.Time = utils.Lerp(0, p.ShutterOpen, p.ShutterClose)

	if screenX == float64(DEBUG_X) && screenY == float64(DEBUG_Y) {
		fmt.Printf("Before transformation ray:\n%v\n", ray)
	}

	return p.CameraToWorld.Ray(ray), 1.0
}

func NewPerspectiveCamera(
	camToWorld *transform.Transform,
	screenWindow [4]float64,
	sopen, sclose, lensr, focald, fov float64,
	f film.Film) *PerspectiveCamera {

	out := &PerspectiveCamera{}

	cam := NewProjectiveCamera(camToWorld,
		transform.Perspective(fov, 1e-2, 100.0),
		screenWindow, sopen, sclose, lensr, focald, f)

	out.CameraToWorld = cam.CameraToWorld
	out.ShutterOpen = cam.ShutterOpen
	out.ShutterClose = cam.ShutterClose
	out.Film = cam.Film
	out.cameraToScreen = cam.cameraToScreen
	out.rasterToCamera = cam.rasterToCamera
	out.screenToRaster = cam.screenToRaster
	out.rasterToScreen = cam.rasterToScreen
	out.lensRadius = cam.lensRadius
	out.focalDistance = cam.focalDistance

	out.dxCamera = out.rasterToCamera.Point(geometry.NewPoint(1, 0, 0)).Minus(
		out.rasterToCamera.Point(geometry.NewPoint(0, 0, 0)))

	out.dyCamera = out.rasterToCamera.Point(geometry.NewPoint(0, 1, 0)).Minus(
		out.rasterToCamera.Point(geometry.NewPoint(0, 0, 0)))

	return out
}
