package shape

import (
	"fmt"
	"math"

	"github.com/ironsmile/raytracer/geometry"
)

func IntersectMultiple(objects []Shape, ray geometry.Ray) (
	prim Shape, retdist float64, normal geometry.Vector) {

	retdist = math.MaxFloat64

	for sInd, shape := range objects {

		bbox := shape.GetObjectBBox()

		if bbox != nil {
			intersected, tNear, _ := bbox.IntersectP(ray)
			if !intersected || tNear > retdist {
				continue
			}
		}

		if shape == nil {
			fmt.Printf("Shape with index %d was nil\n", sInd)
			continue
		}

		res, resDist, resNormal := shape.Intersect(ray, retdist)

		if res == HIT && resDist < retdist {
			prim = shape
			retdist = resDist
			normal = resNormal
		}
	}

	return
}
