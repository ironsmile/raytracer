package shape

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
)

// Quad represents a convex quadrilateral in 3D space. This type implements the Shape
// interface
type Quad struct {
	BasicShape

	vertices [4]geometry.Vector
}

// NewQuad returns a initialized quad structure
func NewQuad(v1, v2, v3, v4 geometry.Vector) *Quad {
	q := &Quad{
		vertices: [4]geometry.Vector{v1, v2, v3, v4},
	}

	q.bbox = bbox.FromPoint(v1)
	q.bbox = bbox.UnionPoint(q.bbox, v2)
	q.bbox = bbox.UnionPoint(q.bbox, v3)
	q.bbox = bbox.UnionPoint(q.bbox, v4)

	return q
}

// Intersect implements the Shape interface for quad face in 3D space. It is based on the
// Ares Lagae and Philip Dutre (2005) algorithm.
func (q *Quad) Intersect(ray geometry.Ray, dg *DifferentialGeometry) bool {
	e01 := q.vertices[1].Minus(q.vertices[0])
	e03 := q.vertices[3].Minus(q.vertices[0])
	p := ray.Direction.Cross(e03)
	det := e01.Dot(p)
	if det == 0 {
		return false
	}
	invDet := 1 / det
	t := ray.Origin.Minus(q.vertices[0])
	alfa := t.Dot(p) * invDet
	if alfa < 0 || alfa > 1 {
		return false
	}
	w := t.Cross(e01)
	beta := ray.Direction.Dot(w) * invDet
	if beta < 0 || beta > 1 {
		return false
	}

	if alfa+beta > 1 {
		e21 := q.vertices[1].Minus(q.vertices[2])
		e23 := q.vertices[3].Minus(q.vertices[2])

		pp := ray.Direction.Cross(e21)
		detp := e23.Dot(pp)
		if detp == 0 {
			return false
		}
		invDetp := 1 / detp
		tp := ray.Origin.Minus(q.vertices[2])
		alfap := tp.Dot(pp) * invDetp
		if alfap < 0 {
			return false
		}
		qp := tp.Cross(e23)
		betap := ray.Direction.Dot(qp) * invDetp
		if betap < 0 {
			return false
		}
	}

	tDist := e03.Dot(w) * invDet

	if tDist < ray.Mint || tDist > ray.Maxt {
		return false
	}

	if dg == nil {
		return true
	}

	dg.Shape = q
	dg.Distance = tDist

	return true
}

// NormalAt implements the Shape interface
func (q *Quad) NormalAt(geometry.Vector) geometry.Vector {
	e01 := q.vertices[1].Minus(q.vertices[0])
	e03 := q.vertices[3].Minus(q.vertices[0])
	return e01.Cross(e03).Normalize()
}

// IntersectP implements the Shape interface
func (q *Quad) IntersectP(ray geometry.Ray) bool {
	return q.Intersect(ray, nil)
}
