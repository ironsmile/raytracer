package scene

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
)

type PlanePrim struct {
	BasePrimitive

	Normal   *geometry.Vector
	Distance float64
}

func (p *PlanePrim) GetType() int {
	return PLANE
}

func NewPlanePrim(normal *geometry.Vector, d float64) *PlanePrim {
	return &PlanePrim{Normal: normal.NormalizeIP(), Distance: d}
}

func (p *PlanePrim) GetNormal(_ *geometry.Point) *geometry.Vector {
	if p.Mat.Refl > 0.0 {
		return p.Normal.Neg()
	}
	return p.Normal
}

func (p *PlanePrim) GetDistance() float64 {
	return p.Distance
}

func (p *PlanePrim) Intersect(ray *geometry.Ray, dist float64) (int, float64, *geometry.Vector) {
	cos := p.Normal.Product(ray.Direction)

	if cos >= 0 {
		return MISS, dist, nil
	}

	dst := -(p.Normal.ProductPoint(ray.Origin) + p.Distance) / cos

	if dst > 0 && dst < dist {
		return HIT, dst, p.GetNormal(nil)
	}

	return MISS, dist, nil
}

func (p *PlanePrim) String() string {
	return fmt.Sprintf("Plane<%s>", p.Name)
}
