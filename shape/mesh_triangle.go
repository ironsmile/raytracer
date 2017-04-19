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
		return HIT, t, m.interpolatedNormal(b1, b2)
	}

	return HIT, t, e1.Cross(e2).Neg().Normalize()
}

func (m *MeshTriangle) interpolatedNormal(b1, b2 float64) geometry.Vector {
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
	return nv1.MultiplyScalar(1 - b1 - b2).Plus(nv2.MultiplyScalar(b1).Plus(nv3.MultiplyScalar(b2)))
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
