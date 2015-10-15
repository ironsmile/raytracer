package scene

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
)

func IntersectPrimitives(primitives []Primitive, ray *geometry.Ray) (Primitive, float64) {
	retdist := 1000000.0
	var prim Primitive = nil

	for sInd, pr := range primitives {

		if pr == nil {
			fmt.Errorf("Primitive with index %d was nil\n", sInd)
		}

		res, resDist := pr.Intersect(ray, retdist)

		if res == HIT && resDist < retdist {
			prim = pr
			retdist = resDist
		}
	}

	return prim, retdist
}
