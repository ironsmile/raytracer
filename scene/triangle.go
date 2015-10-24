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
	triangle.Normal = triangle.edge1.Cross(triangle.edge2).NegIP()

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

	P := ray.Direction.Cross(t.edge2)
	det := t.edge1.Product(P)

	// Not culling:
	if det > -geometry.EPSILON && det < geometry.EPSILON {
		return MISS, dist, nil
	}

	inv_det := 1.0 / det

	T := ray.Origin.Minus(t.Vertices[0])
	u := T.Product(P) * inv_det

	if u < 0.0 || u > 1.0 {
		return MISS, dist, nil
	}

	Q := T.CrossIP(t.edge1)

	v := ray.Direction.Product(Q) * inv_det

	if v < 0.0 || u+v > 1.0 {
		return MISS, dist, nil
	}

	tt := t.edge2.Product(Q) * inv_det

	if tt > geometry.EPSILON && tt <= dist {
		return HIT, tt, t.Normal.Copy()
	}

	return MISS, dist, nil
}
