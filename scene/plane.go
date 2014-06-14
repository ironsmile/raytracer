package scene

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
)

type PlanePrim struct {
	BasePrimitive

	Plane *geometry.Plane
}

func (p *PlanePrim) GetType() int {
	return PLANE
}

func NewPlanePrim(normal geometry.Vector, d float64) *PlanePrim {
	plPrim := new(PlanePrim)
	plPrim.Plane = geometry.NewPlane(normal, d)
	return plPrim
}

func (p *PlanePrim) GetNormal(_ *geometry.Point) *geometry.Vector {
	return p.Plane.N.Copy()
}

func (p *PlanePrim) GetD() float64 {
	return p.Plane.D
}

func (p *PlanePrim) Intersect(ray *geometry.Ray, dist float64) (int, float64) {
	d := p.Plane.N.Product(ray.Direction)

	if d == 0 {
		return MISS, dist
	}

	dst := -(p.Plane.N.ProductPoint(ray.Origin) + p.Plane.D) / d

	if dst > 0 && dst < dist {
		return HIT, dst
	}

	return MISS, dist
}

func (p *PlanePrim) String() string {
	return fmt.Sprintf("Plane<%s>", p.Name)
}
