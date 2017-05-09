package primitive

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

// Primitive is the type which marries the shape to its material. It is resposible for
// the geometry and shading of objects.
type Primitive interface {
	Intersect(geometry.Ray, *Intersection) bool
	IntersectP(geometry.Ray) bool
	IntersectBBoxEdge(geometry.Ray) bool
	GetWorldBBox() *bbox.BBox
	SetTransform(*transform.Transform)
	GetTransforms() (o2w, w2o *transform.Transform)
	CanIntersect() bool
	Refine() []Primitive
	IsLight() bool
	GetLightSource() geometry.Vector
	Shape() shape.Shape
	GetID() uint64
}

// Intersection holds information about a ray–primitive intersection, in-
// cluding information about the differential geometry of the point on the surface, a pointer
// to the Primitive that the ray hit
type Intersection struct {
	DfGeometry shape.DifferentialGeometry
	Primitive  Primitive
}
