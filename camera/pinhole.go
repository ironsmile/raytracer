package camera

import (
	"sync"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/transform"
)

// PinholeCamera is the most basic type of camera. One in which the scene is projected on a
// rectangle and the viewer is a single point behind the screen.
type PinholeCamera struct {
	camToWorld *transform.Transform
	distance   float64
	screen     [4]float64

	rasterW, rasterH float64

	origin geometry.Vector
	lookAt geometry.Vector
	up     geometry.Vector

	sync.RWMutex
}

// GenerateRay creates a ray from the camera source through one single point of the screen
func (p *PinholeCamera) GenerateRay(x, y float64) geometry.Ray {
	posX := p.screen[0] + (x/p.rasterW)*p.screen[1]*2
	posY := p.screen[3] + (y/p.rasterH)*p.screen[2]*2
	ray := geometry.NewRay(
		geometry.NewVector(0, 0, 0),
		geometry.NewVector(posX, posY, p.distance).Normalize(),
	)

	return p.camToWorld.Ray(ray)
}

// Forward moves the camera in its lookAt direction
func (p *PinholeCamera) Forward(speed float64) error {
	p.Lock()
	defer p.Unlock()

	dir := p.lookAt.Minus(p.origin).Normalize().MultiplyScalar(speed)
	p.move(dir)
	return nil
}

// Backward moves the camera opposite its lookAt direction
func (p *PinholeCamera) Backward(speed float64) error {
	p.Lock()
	defer p.Unlock()

	dir := p.lookAt.Minus(p.origin).Normalize().MultiplyScalar(speed).Neg()
	p.move(dir)
	return nil
}

// Left moves the camera to the left relative to its lookAt and up directions
func (p *PinholeCamera) Left(speed float64) error {
	p.Lock()
	defer p.Unlock()

	dir := p.lookAt.Minus(p.origin).Normalize()
	dir = p.up.Cross(dir).MultiplyScalar(speed).Neg()
	p.move(dir)
	return nil
}

// Up moves the camera straight up regardless of its orientation.
func (p *PinholeCamera) Up(speed float64) error {
	p.Lock()
	defer p.Unlock()

	p.move(geometry.NewVector(0, 1, 0).MultiplyScalar(speed))
	return nil
}

// Down moves the camera straight down regardless of its orientation.
func (p *PinholeCamera) Down(speed float64) error {
	p.Lock()
	defer p.Unlock()

	p.move(geometry.NewVector(0, -1, 0).MultiplyScalar(speed))
	return nil
}

// Right moves the camera to the right relative to its lookAt and up directions
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

// Yaw rotats the camera around its up axis
func (p *PinholeCamera) Yaw(angle float64) error {
	p.Lock()
	defer p.Unlock()

	p.rotate(transform.RotateY(angle))
	return nil
}

// Pitch rotates the camera around on the lookAt axis
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

// NewPinhole returns a new camera which is set up for writing in particular output
func NewPinhole(
	camPosition geometry.Vector,
	camLookAtPoint geometry.Vector,
	camUp geometry.Vector,
	dist float64,
	width float64,
	height float64,
) *PinholeCamera {
	cam := &PinholeCamera{
		origin:   camPosition,
		lookAt:   camLookAtPoint,
		up:       camUp,
		distance: dist}

	cam.computeMatrix()

	cam.rasterW = width
	cam.rasterH = height

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
