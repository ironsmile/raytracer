package camera

import (
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/transform"
)

type PinholeCamera struct {
	Film       film.Film
	camToWorld *transform.Transform
	distance   float64
	screen     [4]float64

	rasterW, rasterH float64
}

func (p *PinholeCamera) GenerateRay(screenX, screenY float64) (*geometry.Ray, float64) {
	posX := p.screen[0] + (screenX/p.rasterW)*p.screen[1]*2
	posY := p.screen[3] + (screenY/p.rasterH)*p.screen[2]*2

	origin := geometry.NewPoint(0, 0, 0)
	dir := geometry.Normalize(geometry.NewVector(posX, posY, p.distance))
	ray := geometry.NewRay(*origin, *dir)

	return p.camToWorld.Ray(ray), 1.0
}

func NewPinholeCamera(camToWorld *transform.Transform, dist float64,
	f film.Film) *PinholeCamera {
	cam := &PinholeCamera{Film: f, camToWorld: camToWorld, distance: dist}

	cam.rasterW = float64(f.Width())
	cam.rasterH = float64(f.Height())

	frame := cam.rasterW / cam.rasterH

	if frame > 1.0 {
		cam.screen[0] = -frame
		cam.screen[1] = frame
		cam.screen[2] = -1.0
		cam.screen[3] = 1.0
	} else {
		cam.screen[0] = -1.0
		cam.screen[1] = 1.0
		cam.screen[2] = -1.0 / frame
		cam.screen[3] = 1.0 / frame
	}

	return cam
}
