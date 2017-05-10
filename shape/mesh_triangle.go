package shape

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"

	"github.com/mokiat/go-data-front/decoder/obj"
)

// MeshTriangle is a triangle defined in a object mesh
type MeshTriangle struct {
	BasicShape

	mesh *Mesh
	face *obj.Face
}

// Intersect implements the Shape interface
func (m *MeshTriangle) Intersect(ray geometry.Ray, dg *DifferentialGeometry) bool {

	p1, p2, p3 := m.getPoints()
	e1 := p2.Minus(p1)
	e2 := p3.Minus(p1)

	s1 := ray.Direction.Cross(e2)
	divisor := s1.Product(e1)

	if divisor == 0 {
		return false
	}

	invDivisor := 1 / divisor

	d := ray.Origin.Minus(p1)
	b1 := d.Product(s1) * invDivisor

	if b1 < 0 || b1 > 1 {
		return false
	}

	s2 := d.Cross(e1)
	b2 := ray.Direction.Product(s2) * invDivisor

	if b2 < 0 || b1+b2 > 1 {
		return false
	}

	t := e2.Product(s2) * invDivisor

	if t < ray.Mint || t > ray.Maxt {
		return false
	}

	if dg == nil {
		return true
	}

	dg.Shape = m
	dg.Distance = t

	return true
}

func (m *MeshTriangle) barycentric(p geometry.Vector) (u, v, w float64) {
	p1, p2, p3 := m.getPoints()

	v0 := p2.Minus(p1)
	v1 := p3.Minus(p1)
	v2 := p.Minus(p1)
	d00 := v0.Dot(v0)
	d01 := v0.Dot(v1)
	d11 := v1.Dot(v1)
	d20 := v2.Dot(v0)
	d21 := v2.Dot(v1)
	d := d00*d11 - d01*d01
	v = (d11*d20 - d01*d21) / d
	w = (d00*d21 - d01*d20) / d
	u = 1 - v - w
	return
}

// NormalAt implements the Shape interface
func (m *MeshTriangle) NormalAt(p geometry.Vector) geometry.Vector {

	if m.face.References[0].HasNormal() && m.face.References[1].HasNormal() &&
		m.face.References[2].HasNormal() {

		u, v, w := m.barycentric(p)
		return m.interpolatedNormal(u, v, w)
	}

	p1, p2, p3 := m.getPoints()
	e1 := p2.Minus(p1)
	e2 := p3.Minus(p1)

	return e1.Cross(e2).Neg().Normalize()
}

func (m *MeshTriangle) interpolatedNormal(u, v, w float64) geometry.Vector {
	// Phong interpolation
	// http://paulbourke.net/texture_colour/interpolation/

	n1 := m.mesh.model.GetNormalFromReference(m.face.References[0])
	n2 := m.mesh.model.GetNormalFromReference(m.face.References[1])
	n3 := m.mesh.model.GetNormalFromReference(m.face.References[2])

	nv1 := geometry.NewVector(n1.X, n1.Y, n1.Z).Normalize()
	nv2 := geometry.NewVector(n2.X, n2.Y, n2.Z).Normalize()
	nv3 := geometry.NewVector(n3.X, n3.Y, n3.Z).Normalize()
	return nv1.MultiplyScalar(u).Plus(nv2.MultiplyScalar(v).Plus(nv3.MultiplyScalar(w)))
}

// IntersectP implements the Shape interface
func (m *MeshTriangle) IntersectP(ray geometry.Ray) bool {
	return m.Intersect(ray, nil)
}

func (m *MeshTriangle) getPoints() (p1, p2, p3 geometry.Vector) {
	v1 := m.mesh.model.GetVertexFromReference(m.face.References[0])
	p1.X, p1.Y, p1.Z = v1.X, v1.Y, v1.Z

	v2 := m.mesh.model.GetVertexFromReference(m.face.References[1])
	p2.X, p2.Y, p2.Z = v2.X, v2.Y, v2.Z

	v3 := m.mesh.model.GetVertexFromReference(m.face.References[2])
	p3.X, p3.Y, p3.Z = v3.X, v3.Y, v3.Z

	return
}

// MaterialAt implements the Shape interface
func (m *MeshTriangle) MaterialAt(p geometry.Vector) *mat.Material {
	return m.mesh.MaterialAt(p)
}

// SetMaterial implements Shape interface
func (m *MeshTriangle) SetMaterial(mtl mat.Material) {
	m.mesh.SetMaterial(mtl)
}

// NewMeshTriangle returns a newly created mesh triangle
func NewMeshTriangle(m *Mesh, face *obj.Face) *MeshTriangle {
	mt := MeshTriangle{
		mesh: m,
		face: face,
	}

	p1, p2, p3 := mt.getPoints()

	mt.bbox = bbox.FromPoint(p1)
	mt.bbox = bbox.UnionPoint(mt.bbox, p2)
	mt.bbox = bbox.UnionPoint(mt.bbox, p3)

	return &mt
}
