package scene

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
)

type Triangle struct {
	BasePrimitive

	edge1 *geometry.Vector
	edge2 *geometry.Vector

	Normal   *geometry.Vector
	Vertices [3]*geometry.Point
}

func NewTriangle(vertices [3]*geometry.Point) *Triangle {
	triangle := &Triangle{Vertices: vertices}

	triangle.edge1 = vertices[1].Minus(vertices[0])
	triangle.edge2 = vertices[2].Minus(vertices[0])
	triangle.Normal = triangle.edge1.Cross(triangle.edge2).NegIP().NormalizeIP()

	return triangle
}

func (t *Triangle) GetType() int {
	return TRIANGLE
}

func (t *Triangle) GetNormal(_ *geometry.Point) *geometry.Vector {
	return t.Normal
}

func (t *Triangle) String() string {
	return fmt.Sprintf("Triangle<%s>: %+v", t.Name, t.Vertices)
}

func (t *Triangle) Intersect(ray *geometry.Ray, dist float64) (int, float64, *geometry.Vector) {
	// Implements Möller–Trumbore ray-triangle intersection algorithm:
	// https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm

	s1 := ray.Direction.Cross(t.edge2)
	divisor := t.edge1.Product(s1)

	// Not culling:
	if divisor > -geometry.EPSILON && divisor < geometry.EPSILON {
		return MISS, dist, nil
	}

	invDivisor := 1.0 / divisor

	s := ray.Origin.Minus(t.Vertices[0])
	b1 := s.Product(s1) * invDivisor

	if b1 < 0.0 || b1 > 1.0 {
		return MISS, dist, nil
	}

	s2 := s.CrossIP(t.edge1)
	b2 := ray.Direction.Product(s2) * invDivisor

	if b2 < 0.0 || b1+b2 > 1.0 {
		return MISS, dist, nil
	}

	tt := t.edge2.Product(s2) * invDivisor

	if tt < geometry.EPSILON || tt > dist {
		return MISS, dist, nil
	}

	return HIT, tt, t.Normal.Copy()
}
