package shape

import "github.com/ironsmile/raytracer/geometry"

type Plane struct {
	BasicShape

	Normal   geometry.Vector
	Distance float64
}

func NewPlane(normal geometry.Vector, d float64) *Plane {
	return &Plane{Normal: normal.Normalize(), Distance: d}
}

func (p *Plane) GetNormal(_ *geometry.Vector) geometry.Vector {
	return p.Normal
}

func (p *Plane) GetDistance() float64 {
	return p.Distance
}

func (p *Plane) Intersect(ray geometry.Ray, dist float64) (int, float64, geometry.Vector) {
	cos := p.Normal.Product(ray.Direction)

	if cos >= 0 {
		return MISS, dist, ray.Direction
	}

	dst := -(p.Normal.Product(ray.Origin) + p.Distance) / cos

	if dst > 0 && dst < dist {
		return HIT, dst, p.GetNormal(nil)
	}

	return MISS, dist, ray.Direction
}
