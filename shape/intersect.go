package shape

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
)

func IntersectMultiple(objects []Shape, ray geometry.Ray) (
	prim Shape, retdist float64, normal geometry.Vector) {

	retdist = ray.Maxt

	for sInd, shape := range objects {

		bbox := shape.GetObjectBBox()

		if bbox != nil {
			intersected, tNear, _ := bbox.IntersectP(ray)
			if !intersected || tNear > ray.Maxt {
				continue
			}
		}

		if shape == nil {
			fmt.Printf("Shape with index %d was nil\n", sInd)
			continue
		}

		res, resDist, resNormal := shape.Intersect(ray)

		if res != HIT || resDist > ray.Maxt || resDist < ray.Mint {
			continue
		}

		prim = shape
		retdist = resDist
		ray.Maxt = retdist
		normal = resNormal
	}

	return
}
