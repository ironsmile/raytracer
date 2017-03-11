package primitive

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/shape"
)

type PlanePrim struct {
	BasePrimitive
}

func (p *PlanePrim) GetType() int {
	return PLANE
}

func NewPlanePrim(normal *geometry.Vector, d float64) *PlanePrim {
	planeShape := &shape.PlanePrim{Normal: normal.NormalizeIP(), Distance: d}
	planePrim := &PlanePrim{}
	planePrim.shape = planeShape
	return planePrim
}

func (p *PlanePrim) Intersect(ray *geometry.Ray, dist float64) (int, float64, *geometry.Vector) {
	return p.shape.Intersect(ray, dist)
}

func (p *PlanePrim) String() string {
	return fmt.Sprintf("Plane<%s>", p.Name)
}
