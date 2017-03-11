package shape

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
)

func IntersectMultiple(objects []Shape, ray *geometry.Ray) (
	prim Shape, retdist float64, normal *geometry.Vector) {

	retdist = 1000000.0

	for sInd, pr := range objects {

		if pr == nil {
			fmt.Printf("Shape with index %d was nil\n", sInd)
			continue
		}

		res, resDist, resNormal := pr.Intersect(ray, retdist)

		if res == HIT && resDist < retdist {
			prim = pr
			retdist = resDist
			normal = resNormal
		}
	}

	return
}
