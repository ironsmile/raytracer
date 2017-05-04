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

	dg.Distance = t

	if m.face.References[0].HasNormal() && m.face.References[1].HasNormal() &&
		m.face.References[2].HasNormal() {
		dg.Normal = m.interpolatedNormal(b1, b2)
	} else {
		dg.Normal = e1.Cross(e2).Neg().Normalize()
	}

	return true
}

func (m *MeshTriangle) interpolatedNormal(b1, b2 float64) geometry.Vector {
	// Phong interpolation
	// http://paulbourke.net/texture_colour/interpolation/

	n1 := m.mesh.model.GetNormalFromReference(m.face.References[0])
	n2 := m.mesh.model.GetNormalFromReference(m.face.References[1])
	n3 := m.mesh.model.GetNormalFromReference(m.face.References[2])

	nv1 := geometry.NewVector(n1.X, n1.Y, n1.Z).Normalize()
	nv2 := geometry.NewVector(n2.X, n2.Y, n2.Z).Normalize()
	nv3 := geometry.NewVector(n3.X, n3.Y, n3.Z).Normalize()

	//!TODO: implement proper coefficients for phong shading. Obviously 0.5
	// for all points is a lie
	return nv1.MultiplyScalar(1 - b1 - b2).Plus(nv2.MultiplyScalar(b1).Plus(nv3.MultiplyScalar(b2)))
}

// IntersectP implements the Shape interface
func (m *MeshTriangle) IntersectP(ray geometry.Ray) bool {
	return m.Intersect(ray, nil)
}

func (m *MeshTriangle) getPoints() (p1, p2, p3 geometry.Vector) {
	v1 := m.mesh.model.GetVertexFromReference(m.face.References[0])
	p1 = geometry.NewVector(v1.X, v1.Y, v1.Z)

	v2 := m.mesh.model.GetVertexFromReference(m.face.References[1])
	p2 = geometry.NewVector(v2.X, v2.Y, v2.Z)

	v3 := m.mesh.model.GetVertexFromReference(m.face.References[2])
	p3 = geometry.NewVector(v3.X, v3.Y, v3.Z)

	return
}

// GetMaterial implements the Shape interface
func (m *MeshTriangle) GetMaterial() *mat.Material {
	return m.mesh.GetMaterial()
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
