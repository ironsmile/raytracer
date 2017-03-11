package shape

import "github.com/ironsmile/raytracer/geometry"

type PlanePrim struct {
	Normal   *geometry.Vector
	Distance float64
}

func NewPlanePrim(normal *geometry.Vector, d float64) *PlanePrim {
	return &PlanePrim{Normal: normal.NormalizeIP(), Distance: d}
}

func (p *PlanePrim) GetNormal(_ *geometry.Point) *geometry.Vector {
	// Commented during shapre/primitive division
	// if p.Mat.Refl > 0.0 {
	// 	return p.Normal.Neg()
	// }
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
