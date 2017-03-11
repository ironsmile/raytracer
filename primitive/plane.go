package primitive

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

type Plane struct {
	BasePrimitive
}

func (p *Plane) GetType() int {
	return PLANE
}

func NewPlane(normal *geometry.Vector, d float64) *Plane {
	planePrim := &Plane{}
	planePrim.shape = shape.NewPlane(normal, d)
	planePrim.SetTransform(transform.Identity())
	return planePrim
}

func (p *Plane) String() string {
	return fmt.Sprintf("Plane<%s>", p.Name)
}
