package primitive

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/shape"
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
	return planePrim
}

func (p *Plane) Intersect(ray *geometry.Ray, dist float64) (int, float64, *geometry.Vector) {
	return p.shape.Intersect(ray, dist)
}

func (p *Plane) String() string {
	return fmt.Sprintf("Plane<%s>", p.Name)
}
