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

	origin geometry.Vector
	lookAt geometry.Vector
	up     geometry.Vector

	sync.RWMutex
}

func (p *PinholeCamera) GenerateRay(x, y float64) (geometry.Ray, float64) {
	p.RLock()
	defer p.RUnlock()

	posX := p.screen[0] + (x/p.rasterW)*p.screen[1]*2
	posY := p.screen[3] + (y/p.rasterH)*p.screen[2]*2
	ray := geometry.Ray{
		Origin:    geometry.NewVector(0, 0, 0),
		Direction: geometry.NewVector(posX, posY, p.distance).Normalize(),
	}

	return p.camToWorld.Ray(ray), 1.0
}

func (p *PinholeCamera) Forward(speed float64) error {
	p.Lock()
	defer p.Unlock()

	dir := p.lookAt.Minus(p.origin).Normalize().MultiplyScalar(speed)
	p.move(dir)
	return nil
}

func (p *PinholeCamera) Backward(speed float64) error {
	p.Lock()
	defer p.Unlock()

	dir := p.lookAt.Minus(p.origin).Normalize().MultiplyScalar(speed).Neg()
	p.move(dir)
	return nil
}

func (p *PinholeCamera) Left(speed float64) error {
	p.Lock()
	defer p.Unlock()

	dir := p.lookAt.Minus(p.origin).Normalize()
	dir = p.up.Cross(dir).MultiplyScalar(speed).Neg()
	p.move(dir)
	return nil
}

func (p *PinholeCamera) Right(speed float64) error {
	p.Lock()
	defer p.Unlock()

	dir := p.lookAt.Minus(p.origin).Normalize()
	dir = p.up.Cross(dir).MultiplyScalar(speed)
	p.move(dir)
	return nil
}

func (p *PinholeCamera) move(dir geometry.Vector) {
	p.origin = p.origin.Plus(dir)
	p.lookAt = p.lookAt.Plus(dir)
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
	p.lookAt = p.camToWorld.Point(rotation.Point(p.camToWorld.Inverse().Point(p.lookAt)))
	p.computeMatrix()
}

func NewPinhole(camPosition, camLookAtPoint geometry.Vector,
	camUp geometry.Vector, dist float64, f film.Film) *PinholeCamera {

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
