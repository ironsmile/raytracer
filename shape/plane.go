package shape

import "github.com/ironsmile/raytracer/geometry"

type Plane struct {
	Normal   *geometry.Vector
	Distance float64
}

func NewPlane(normal *geometry.Vector, d float64) *Plane {
	return &Plane{Normal: normal.NormalizeIP(), Distance: d}
}

func (p *Plane) GetNormal(_ *geometry.Point) *geometry.Vector {
	// Commented during shapre/primitive division
	// if p.Mat.Refl > 0.0 {
	// 	return p.Normal.Neg()
	// }
	return p.Normal
}

func (p *Plane) GetDistance() float64 {
	return p.Distance
}

func (p *Plane) Intersect(ray *geometry.Ray, dist float64) (int, float64, *geometry.Vector) {
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
