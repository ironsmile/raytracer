package shape

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/momchil-atanasov/go-data-front/decoder/obj"

	"github.com/ironsmile/raytracer/geometry"
)

// MeshTriangle is a triangle defined in a object mesh
type MeshTriangle struct {
	BasicShape

	model *obj.Model
	face  *obj.Face
}

// Intersect implements the Shape interface
func (m *MeshTriangle) Intersect(ray geometry.Ray) (int, float64, geometry.Vector) {

	p1, p2, p3 := m.getPoints()
	e1 := p2.Minus(p1)
	e2 := p3.Minus(p1)

	s1 := ray.Direction.Cross(e2)
	divisor := s1.Product(e1)

	if divisor == 0 {
		return MISS, 0, ray.Direction
	}

	invDivisor := 1 / divisor

	d := ray.Origin.Minus(p1)
	b1 := d.Product(s1) * invDivisor

	if b1 < 0 || b1 > 1 {
		return MISS, 0, ray.Direction
	}

	s2 := d.Cross(e1)
	b2 := ray.Direction.Product(s2) * invDivisor

	if b2 < 0 || b1+b2 > 1 {
		return MISS, 0, ray.Direction
	}

	t := e2.Product(s2) * invDivisor

	if t < ray.Mint || t > ray.Maxt {
		return MISS, 0, ray.Direction
	}

	if m.face.References[0].HasNormal() && m.face.References[1].HasNormal() &&
		m.face.References[2].HasNormal() {
		return HIT, t, m.interpolatedNormal(ray.At(t))
	}

	if !m.face.References[0].HasTexCoord() {
		return HIT, t, e1.Cross(e2).Neg().Normalize()
	}

	uv0 := m.model.GetTexCoordFromReference(m.face.References[0])
	uv1 := m.model.GetTexCoordFromReference(m.face.References[1])
	uv2 := m.model.GetTexCoordFromReference(m.face.References[2])

	du1 := uv0.U - uv2.U
	du2 := uv1.U - uv2.U
	dv1 := uv0.V - uv2.V
	dv2 := uv1.V - uv2.V

	dp1 := p1.Minus(p3)
	dp2 := p2.Minus(p3)

	var dpdu, dpdv geometry.Vector
	determinant := du1 * dv2 * dv1 * du2

	if determinant == 0 {
		dpdu, dpdv = geometry.CoordinateSystem(e2.Cross(e1).Normalize())
	} else {
		invdet := 1 / determinant
		dpdu = dp1.MultiplyScalar(dv2).Minus(
			dp2.MultiplyScalar(dv1)).MultiplyScalar(invdet)
		dpdv = dp1.MultiplyScalar(-du2).Plus(
			dp2.MultiplyScalar(du1)).MultiplyScalar(invdet)
	}

	return HIT, t, dpdu.Cross(dpdv).Normalize()
}

func (m *MeshTriangle) interpolatedNormal(pi geometry.Vector) geometry.Vector {
	// Phong interpolation
	// http://paulbourke.net/texture_colour/interpolation/

	n1 := m.model.GetNormalFromReference(m.face.References[0])
	n2 := m.model.GetNormalFromReference(m.face.References[1])
	n3 := m.model.GetNormalFromReference(m.face.References[2])

	nv1 := geometry.NewVector(n1.X, n1.Y, n1.Z).Normalize()
	nv2 := geometry.NewVector(n2.X, n2.Y, n2.Z).Normalize()
	nv3 := geometry.NewVector(n3.X, n3.Y, n3.Z).Normalize()

	//!TODO: implement proper coefficients for phong shading. Obviously 0.5
	// for all points is a lie
	nv12 := geometry.Lerp(nv1, nv2, 0.5).Normalize()
	nv13 := geometry.Lerp(nv1, nv3, 0.5).Normalize()

	return geometry.Lerp(nv12, nv13, 0.5).Normalize()
}

// IntersectP implements the Shape interface
func (m *MeshTriangle) IntersectP(ray geometry.Ray) bool {
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

	return true
}

func (m *MeshTriangle) getPoints() (p1, p2, p3 geometry.Vector) {
	v1 := m.model.GetVertexFromReference(m.face.References[0])
	p1 = geometry.NewVector(v1.X, v1.Y, v1.Z)

	v2 := m.model.GetVertexFromReference(m.face.References[1])
	p2 = geometry.NewVector(v2.X, v2.Y, v2.Z)

	v3 := m.model.GetVertexFromReference(m.face.References[2])
	p3 = geometry.NewVector(v3.X, v3.Y, v3.Z)

	return
}

// NewMeshTriangle returns a newly created mesh triangle
func NewMeshTriangle(m *obj.Model, face *obj.Face) *MeshTriangle {
	mt := MeshTriangle{
		model: m,
		face:  face,
	}

	p1, p2, p3 := mt.getPoints()

	mt.bbox = bbox.FromPoint(p1)
	mt.bbox = bbox.UnionPoint(mt.bbox, p2)
	mt.bbox = bbox.UnionPoint(mt.bbox, p3)

	return &mt
}
