package shape

import "github.com/ironsmile/raytracer/geometry"

type Triangle struct {
	BasicShape

	edge1 geometry.Vector
	edge2 geometry.Vector

	Normal   geometry.Vector
	Vertices [3]geometry.Vector
}

func NewTriangle(vertices [3]geometry.Vector) *Triangle {
	triangle := &Triangle{Vertices: vertices}

	triangle.edge1 = vertices[1].Minus(vertices[0])
	triangle.edge2 = vertices[2].Minus(vertices[0])
	triangle.Normal = triangle.edge1.Cross(triangle.edge2).Neg().Normalize()

	return triangle
}

func (t *Triangle) Intersect(ray geometry.Ray, dist float64) (int, float64, geometry.Vector) {
	// Implements Möller–Trumbore ray-triangle intersection algorithm:
	// https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm
	var outNormal geometry.Vector

	s1 := ray.Direction.Cross(t.edge2)
	divisor := t.edge1.Product(s1)

	// Not culling:
	if divisor > -geometry.EPSILON && divisor < geometry.EPSILON {
		return MISS, dist, outNormal
	}

	invDivisor := 1.0 / divisor

	s := ray.Origin.Minus(t.Vertices[0])
	b1 := s.Product(s1) * invDivisor

	if b1 < 0.0 || b1 > 1.0 {
		return MISS, dist, outNormal
	}

	s2 := s.Cross(t.edge1)
	b2 := ray.Direction.Product(s2) * invDivisor

	if b2 < 0.0 || b1+b2 > 1.0 {
		return MISS, dist, outNormal
	}

	tt := t.edge2.Product(s2) * invDivisor

	if tt < geometry.EPSILON || tt > dist {
		return MISS, dist, outNormal
	}

	return HIT, tt, t.Normal
}
