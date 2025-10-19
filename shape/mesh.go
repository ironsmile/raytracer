package shape

import (
	"fmt"

	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"

	"github.com/mokiat/go-data-front/decoder/obj"
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
		switch len(face.References) {
		case 3:
			m.bbox = bbox.Union(m.bbox, NewMeshTriangle(&m, face).GetObjectBBox())
		case 4:
			m.bbox = bbox.Union(m.bbox, NewMeshQuad(&m, face).GetObjectBBox())
		}
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
	meshFaces := make([]Shape, 0, len(m.mesh.Faces))
	for faceIndex, face := range m.mesh.Faces {
		refs := len(face.References)
		if refs < 3 {
			panic(fmt.Sprintf(
				"face %d [mesh: %+v] has %d points, cannot load it",
				faceIndex, m.mesh.MaterialName, len(face.References)))
		} else if refs == 3 {
			meshFaces = append(meshFaces, NewMeshTriangle(m, face))
			// } else if refs == 4 {
			// 	meshFaces = append(meshFaces, NewMeshQuad(m, face))
		} else {
			for i := 1; i < refs-1; i++ {
				trigFace := obj.Face{
					References: []obj.Reference{
						face.References[0],
						face.References[i],
						face.References[i+1],
					},
				}
				meshFaces = append(meshFaces, NewMeshTriangle(m, &trigFace))
			}
		}
	}
	return meshFaces
}
