package scene

import (
	"fmt"

	"github.com/ironsmile/raytracer/common"
)

type PlanePrim struct {
	BasePrimitive

	Plane *common.Plane
}

func (p *PlanePrim) GetType() int {
	return PLANE
}

func NewPlanePrim(normal common.Vector, d float64) *PlanePrim {
	plPrim := new(PlanePrim)
	plPrim.Plane = common.NewPlane(normal, d)
	return plPrim
}

func (p *PlanePrim) GetNormal(_ *common.Vector) *common.Vector {
	return p.Plane.N.Copy()
}

func (p *PlanePrim) GetD() float64 {
	return p.Plane.D
}

func (p *PlanePrim) Intersect(ray *common.Ray, dist float64) (int, float64) {
	d := p.Plane.N.Product(ray.Direction)

	if d == 0 {
		return MISS, dist
	}

	dst := -(p.Plane.N.Product(ray.Origin) + p.Plane.D) / d

	if dst > 0 {
		if dst < dist {
			return HIT, dst
		}
	}

	return MISS, dist
}

func (p *PlanePrim) String() string {
	return fmt.Sprintf("Plane<%s>", p.Name)
}
