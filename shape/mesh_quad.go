package shape

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"
	"github.com/mokiat/go-data-front/decoder/obj"
)

// MeshQuad is a quad defined by a mesh face.
type MeshQuad struct {
	BasicShape

	mesh *Mesh
	face *obj.Face
}

// NewMeshQuad returns a MeshQuad for a face in a mesh.
func NewMeshQuad(m *Mesh, face *obj.Face) *MeshQuad {
	mq := MeshQuad{
		mesh: m,
		face: face,
	}

	p1, p2, p3, p4 := mq.getPoints()

	mq.bbox = bbox.FromPoint(p1)
	mq.bbox = bbox.UnionPoint(mq.bbox, p2)
	mq.bbox = bbox.UnionPoint(mq.bbox, p3)
	mq.bbox = bbox.UnionPoint(mq.bbox, p4)

	return &mq
}

// Intersect implements the [Shape] interface.
func (m *MeshQuad) Intersect(ray geometry.Ray, dg *DifferentialGeometry) bool {
	p0, p1, p2, p3 := m.getPoints()

	e01 := p1.Minus(p0)
	e03 := p3.Minus(p0)
	p := ray.Direction.Cross(e03)
	det := e01.Dot(p)
	if det == 0 {
		return false
	}
	invDet := 1 / det
	t := ray.Origin.Minus(p0)
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
		e21 := p1.Minus(p2)
		e23 := p3.Minus(p2)

		pp := ray.Direction.Cross(e21)
		detp := e23.Dot(pp)
		if detp == 0 {
			return false
		}
		invDetp := 1 / detp
		tp := ray.Origin.Minus(p2)
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

	dg.Shape = m
	dg.Distance = tDist

	return true
}

// IntersectP implements the [Shape] interface.
func (m *MeshQuad) IntersectP(ray geometry.Ray) bool {
	return m.Intersect(ray, nil)
}

// NormalAt implements the [Shape] interface.
func (m *MeshQuad) NormalAt(geometry.Vector) geometry.Vector {
	if m.face.References[0].HasNormal() {
		n1 := m.mesh.model.GetNormalFromReference(m.face.References[0])
		return geometry.NewVector(n1.X, n1.Y, n1.Z)
	}

	p0, p1, _, p3 := m.getPoints()

	e01 := p1.Minus(p0)
	e03 := p3.Minus(p0)
	return e01.Cross(e03).Normalize()
}

// MaterialAt implements the [Shape] interface.
func (m *MeshQuad) MaterialAt(p geometry.Vector) *mat.Material {
	return m.mesh.MaterialAt(p)
}

// SetMaterial implements [Shape] interface.
func (m *MeshQuad) SetMaterial(mtl mat.Material) {
	m.mesh.SetMaterial(mtl)
}

func (m *MeshQuad) getPoints() (p1, p2, p3, p4 geometry.Vector) {
	v1 := m.mesh.model.GetVertexFromReference(m.face.References[0])
	p1.X, p1.Y, p1.Z = v1.X, v1.Y, v1.Z

	v2 := m.mesh.model.GetVertexFromReference(m.face.References[1])
	p2.X, p2.Y, p2.Z = v2.X, v2.Y, v2.Z

	v3 := m.mesh.model.GetVertexFromReference(m.face.References[2])
	p3.X, p3.Y, p3.Z = v3.X, v3.Y, v3.Z

	v4 := m.mesh.model.GetVertexFromReference(m.face.References[3])
	p4.X, p4.Y, p4.Z = v4.X, v4.Y, v4.Z

	return
}
