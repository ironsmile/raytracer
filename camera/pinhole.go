package camera

import (
	"sync"

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

	origin *geometry.Point
	lookAt *geometry.Point
	up     *geometry.Vector

	sync.RWMutex
}

func (p *PinholeCamera) GenerateRay(x, y float64) (*geometry.Ray, float64) {
	ray := &geometry.Ray{}
	w := p.GenerateRayIP(x, y, ray)
	return ray, w
}

func (p *PinholeCamera) GenerateRayIP(x, y float64, ray *geometry.Ray) float64 {
	p.RLock()
	defer p.RUnlock()

	posX := p.screen[0] + (x/p.rasterW)*p.screen[1]*2
	posY := p.screen[3] + (y/p.rasterH)*p.screen[2]*2

	ray.Origin = geometry.NewPoint(0, 0, 0)
	ray.Direction = geometry.NewVector(posX, posY, p.distance).NormalizeIP()
	p.camToWorld.RayIP(ray)

	return 1.0
}

func (p *PinholeCamera) Forward(speed float64) error {
	p.Lock()
	defer p.Unlock()

	dir := p.lookAt.Minus(p.origin).NormalizeIP().MultiplyScalarIP(speed)
	p.move(dir)
	return nil
}

func (p *PinholeCamera) Backward(speed float64) error {
	p.Lock()
	defer p.Unlock()

	dir := p.lookAt.Minus(p.origin).NormalizeIP().MultiplyScalarIP(speed).NegIP()
	p.move(dir)
	return nil
}

func (p *PinholeCamera) Left(speed float64) error {
	p.Lock()
	defer p.Unlock()

	dir := p.lookAt.Minus(p.origin).NormalizeIP()
	dir = p.up.Cross(dir).MultiplyScalarIP(speed).NegIP()
	p.move(dir)
	return nil
}

func (p *PinholeCamera) Right(speed float64) error {
	p.Lock()
	defer p.Unlock()

	dir := p.lookAt.Minus(p.origin).NormalizeIP()
	dir = p.up.Cross(dir).MultiplyScalarIP(speed)
	p.move(dir)
	return nil
}

func (p *PinholeCamera) move(dir *geometry.Vector) {
	p.origin.PlusVectorIP(dir)
	p.lookAt.PlusVectorIP(dir)
	p.computeMatrix()
}

func (p *PinholeCamera) computeMatrix() {
	p.camToWorld = transform.LookAt(p.origin, p.lookAt, p.up).Inverse()
}

func (p *PinholeCamera) Yaw(angle float64) error {
	p.Lock()
	defer p.Unlock()

	p.rotate(transform.RotateY(angle))
	return nil
}

func (p *PinholeCamera) Pitch(angle float64) error {
	p.Lock()
	defer p.Unlock()

	p.rotate(transform.RotateX(angle))
	return nil
}

func (p *PinholeCamera) rotate(rotation *transform.Transform) {
	p.camToWorld.PointIP(rotation.PointIP(p.camToWorld.Inverse().PointIP(p.lookAt)))
	p.computeMatrix()
}

func NewPinholeCamera(camPosition, camLookAtPoint *geometry.Point,
	camUp *geometry.Vector, dist float64, f film.Film) *PinholeCamera {

	cam := &PinholeCamera{
		Film:     f,
		origin:   camPosition,
		lookAt:   camLookAtPoint,
		up:       camUp,
		distance: dist}

	cam.computeMatrix()

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
