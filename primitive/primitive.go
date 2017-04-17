package primitive

import (
	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"
	"github.com/ironsmile/raytracer/shape"
)

// Primitive is the type which marries the shape to its material. It is resposible for
// the geometry and shading of objects.
type Primitive interface {
	Intersect(geometry.Ray) (pr Primitive, distance float64, normal geometry.Vector)
	IntersectP(geometry.Ray) bool
	IntersectBBoxEdge(geometry.Ray) bool
	GetWorldBBox() *bbox.BBox
	GetColor() *geometry.Color
	GetMaterial() *mat.Material
	IsLight() bool
	GetLightSource() geometry.Vector
	Shape() shape.Shape
}
