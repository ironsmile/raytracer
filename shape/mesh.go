package shape

import (
	"fmt"

	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/momchil-atanasov/go-data-front/decoder/obj"
)

// Mesh represents a single mesh with triangles
type Mesh struct {
	BasicShape

	mesh  *obj.Mesh
	model *obj.Model
}

// NewMesh returns a mesh defined by the
func NewMesh(model *obj.Model, mesh *obj.Mesh) *Mesh {
	m := Mesh{
		model: model,
		mesh:  mesh,
	}

	for _, face := range m.mesh.Faces {
		m.bbox = bbox.Union(m.bbox, NewMeshTriangle(model, face).GetObjectBBox())
	}

	return &m
}

// Intersect implements the Shape interface
func (m *Mesh) Intersect(geometry.Ray, *DifferentialGeometry) bool {
	panic("Cannot Intersect mesh shape directly")
}

// IntersectP implements the Shape interface
func (m *Mesh) IntersectP(geometry.Ray) bool {
	panic("Cannot IntersectP mesh shape directly")
}

// CanIntersect implements the Shape interface
func (m *Mesh) CanIntersect() bool {
	return false
}

// Refine implements the Shape interface
func (m *Mesh) Refine() []Shape {
	meshTriangles := make([]Shape, 0, len(m.mesh.Faces))
	for faceIndex, face := range m.mesh.Faces {
		if len(face.References) != 3 {
			panic(fmt.Sprintf(
				"face %d [mesh: %+v] has %d points, cannot load it",
				faceIndex, m.mesh, len(face.References)))
		}

		meshTriangles = append(meshTriangles, NewMeshTriangle(m.model, face))
	}
	return meshTriangles
}
