package shape

import "github.com/ironsmile/raytracer/geometry"

// IntersectMultiple returns wether the ray intersects a slice of shapes
func IntersectMultiple(objects []Shape, ray geometry.Ray, dg *DifferentialGeometry) bool {

	var hasHit bool

	for _, shape := range objects {

		if bbox := shape.GetObjectBBox(); bbox != nil {
			intersected, _, _ := bbox.IntersectP(ray)
			if !intersected {
				continue
			}
		}

		if ok := shape.Intersect(ray, dg); ok {
			ray.Maxt = dg.Distance
			hasHit = true
		}
	}

	return hasHit
}

// IntersectPMultiple returns wether the ray intersects a slice of shapes and returns
// true or false. It would be faster than IntersectMultiple because it doesn't have to
// calculate intersection data like
func IntersectPMultiple(objects []Shape, ray geometry.Ray) bool {
	for _, shape := range objects {

		if bbox := shape.GetObjectBBox(); bbox != nil {
			intersected, _, _ := bbox.IntersectP(ray)
			if !intersected {
				continue
			}
		}

		if res := shape.IntersectP(ray); res {
			return true
		}
	}

	return false
}
