package primitive

import "github.com/ironsmile/raytracer/geometry"

// IntersectMultiple returns wether the ray intersects a slice of primitives
func IntersectMultiple(primitives []Primitive, ray geometry.Ray, in *Intersection) bool {
	var hasHit bool
	for _, pr := range primitives {

		if ok := pr.Intersect(ray, in); !ok {
			continue
		}

		hasHit = true
		ray.Maxt = in.DfGeometry.Distance
	}

	return hasHit
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
