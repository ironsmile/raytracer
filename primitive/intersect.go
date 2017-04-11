package primitive

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
)

// IntersectMultiple returns wether the ray intersects a slice of primitives
func IntersectMultiple(primitives []Primitive, ray geometry.Ray) (
	prim Primitive,
	retdist float64,
	normal geometry.Vector,
) {
	retdist = ray.Maxt

	for sInd, pr := range primitives {

		if pr == nil {
			fmt.Printf("primitive with index %d was nil\n", sInd)
			continue
		}

		res, resDist, resNormal := pr.Intersect(ray)

		if res == nil || resDist > ray.Maxt || resDist < ray.Mint {
			continue
		}

		prim = pr
		retdist = resDist
		ray.Maxt = resDist
		normal = resNormal
	}

	return
}
