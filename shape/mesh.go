package shape

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/momchil-atanasov/go-data-front/decoder/obj"
)

// Mesh represents a single mesh with triangles
type Mesh struct {
	BasicShape

	triangles []Shape
	model     *obj.Model
}

// NewMesh returns a mesh defined by the
func NewMesh(model *obj.Model, triangles []Shape) *Mesh {
	m := Mesh{
		model:     model,
		triangles: triangles,
	}

	for _, tr := range triangles {
		m.bbox = bbox.Union(m.bbox, tr.GetObjectBBox())
	}

	return &m
}

// Intersect implements the Shape interface
func (m *Mesh) Intersect(ray geometry.Ray, dg *DifferentialGeometry) bool {
	return IntersectMultiple(m.triangles, ray, dg)
}

// IntersectP implements the Shape interface
func (m *Mesh) IntersectP(ray geometry.Ray) bool {
	return IntersectPMultiple(m.triangles, ray)
}

func (m *Mesh) GetAllShapes() []Shape {
	return m.triangles
}
