package primitive

import "github.com/ironsmile/raytracer/geometry"

// IntersectMultiple returns wether the ray intersects a slice of primitives
func IntersectMultiple(primitives []Primitive, ray geometry.Ray) (
	prim Primitive,
	retdist float64,
	normal geometry.Vector,
) {
	for _, pr := range primitives {

		res, resDist, resNormal := pr.Intersect(ray)

		if res == nil {
			continue
		}

		prim = pr
		retdist = resDist
		ray.Maxt = resDist
		normal = resNormal
	}

	return
}

// IntersectPMultiple returns wether the ray intersects a slice of primitives and returns
// true or false. It would be faster than IntersectMultiple because it doesn't have to
// calculate intersection data like
func IntersectPMultiple(primitives []Primitive, ray geometry.Ray) bool {
	for _, pr := range primitives {
		if intersected := pr.IntersectP(ray); intersected {
			return true
		}
	}
	return false
}
