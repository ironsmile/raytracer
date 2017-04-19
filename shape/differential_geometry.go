package shape

import (
	"github.com/ironsmile/raytracer/geometry"
)

// DifferentialGeometry is a self-contained representation for the geometry
// of a particular point on a surface (typically the point of a ray intersection).
// This abstraction needs to hide the particular type of geometric shape the point lies
// on, supplying enough information about the surface point to allow the shading and
// geometric operations in the rest of pbrt to be implemented generically, without the
// need to distinguish between different shape types such as spheres and triangles.
type DifferentialGeometry struct {

	// The normal of the intersection at this point
	Normal geometry.Vector

	// The distance from the ray origin for this intersection
	Distance float64
}
