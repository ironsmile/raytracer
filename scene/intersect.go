package scene

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
)

func IntersectPrimitives(primitives []Primitive, ray *geometry.Ray) (
	prim Primitive, retdist float64, normal *geometry.Vector) {

	retdist = 1000000.0

	for sInd, pr := range primitives {

		if pr == nil {
			fmt.Errorf("Primitive with index %d was nil\n", sInd)
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
