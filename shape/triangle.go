package shape

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
)

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

	triangle.bbox = bbox.FromPoint(vertices[0])
	triangle.bbox = bbox.UnionPoint(triangle.bbox, vertices[1])
	triangle.bbox = bbox.UnionPoint(triangle.bbox, vertices[2])

	return triangle
}

func (t *Triangle) Intersect(ray geometry.Ray, dg *DifferentialGeometry) bool {
	// Implements Möller–Trumbore ray-triangle intersection algorithm:
	// https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm

	s1 := ray.Direction.Cross(t.edge2)
	divisor := t.edge1.Product(s1)

	// Not culling:
	if divisor > -geometry.EPSILON && divisor < geometry.EPSILON {
		return false
	}

	invDivisor := 1.0 / divisor

	s := ray.Origin.Minus(t.Vertices[0])
	b1 := s.Product(s1) * invDivisor

	if b1 < 0.0 || b1 > 1.0 {
		return false
	}

	s2 := s.Cross(t.edge1)
	b2 := ray.Direction.Product(s2) * invDivisor

	if b2 < 0.0 || b1+b2 > 1.0 {
		return false
	}

	tt := t.edge2.Product(s2) * invDivisor

	if tt > ray.Maxt || tt < ray.Mint {
		return false
	}

	if dg == nil {
		return true
	}

	dg.Distance = tt
	dg.Normal = t.Normal

	return true
}

func (t *Triangle) IntersectP(ray geometry.Ray) bool {
	return t.Intersect(ray, nil)
}
