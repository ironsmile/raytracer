package accel

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/primitive"
	"github.com/ironsmile/raytracer/scene/example"
)

func TestAggregatorsIntersections(t *testing.T) {
	rand.Seed(time.Now().Unix())
	prims, _ := example.GetTeapotScene()
	prims = FullyRefinePrimitives(prims)

	tests := []struct {
		name  string
		accel primitive.Primitive
	}{
		{
			name:  "grid",
			accel: NewGrid(prims),
		},
		{
			name:  "bvh",
			accel: NewBVH(prims, 1),
		},
	}

	for _, test := range tests {
		accel := test.accel
		t.Run(test.name, func(t *testing.T) {
			testIntersectionsWithAggregator(t, prims, accel)
		})
	}
}

func testIntersectionsWithAggregator(
	t *testing.T,
	prims []primitive.Primitive,
	accel primitive.Primitive,
) {
	bboxes := make([]*bbox.BBox, len(prims))
	var bb *bbox.BBox

	for i, pr := range prims {
		bboxes[i] = pr.GetWorldBBox()
		bb = bbox.Union(bb, pr.GetWorldBBox())
	}

	bb.Expand(bb.Max.ByAxis(bb.MaximumExtend()) -
		bb.Min.ByAxis(bb.MaximumExtend()))

	var lastHit geometry.Vector

	for i := 0; i < 40000; i++ {
		// Choos ray origin for testing accelerator
		orig := geometry.NewVector(
			randInRange(bb.Min.X, bb.Max.X),
			randInRange(bb.Min.Y, bb.Max.Y),
			randInRange(bb.Min.Z, bb.Max.Z),
		)
		if rand.Int()%4 == 0 {
			orig = lastHit
		}

		// Choose ray direction for testing accelerator
		dir := uniformRandomSphere()

		if rand.Int()%32 == 0 {
			dir.X, dir.Y = 0, 0
		} else if rand.Int()%32 == 0 {
			dir.X, dir.Z = 0, 0
		} else if rand.Int()%32 == 0 {
			dir.Y, dir.Z = 0, 0
		}

		// Choose ray epsilon for testing accelerator 248
		var eps float64

		if rand.Float64() < 0.25 {
			eps = geometry.EPSILON
		}

		ray := geometry.NewRay(orig, dir.Normalize())
		ray.Mint = eps
		rayAll := ray

		var isectAll, isectAccel primitive.Intersection
		var hitAll, inconsistentBounds bool

		for i, pr := range prims {
			if is, _, _ := bboxes[i].IntersectP(rayAll); is {
				if pr.Intersect(rayAll, &isectAll) {
					hitAll = true
					rayAll.Maxt = isectAll.DfGeometry.Distance
				}
			} else if pr.Intersect(rayAll, &isectAll) {
				// It is possible for an accelerator to report a hit even though intersection
				// between the ray and this primitive's bounding box is not reported. This
				// might be because of a rounding error of float calculations while
				// intersecting the bounding box. Cases like this would be ignored.
				inconsistentBounds = true
			}
		}

		if hitAll {
			pi := ray.At(isectAll.DfGeometry.Distance)
			lastHit = pi.Plus(
				isectAll.DfGeometry.Shape.NormalAt(pi).MultiplyScalar(geometry.EPSILON),
			)
		}

		if inconsistentBounds {
			continue
		}

		hitAccel := accel.Intersect(ray, &isectAccel)

		if hitAccel != hitAll ||
			(hitAll && isectAll.DfGeometry.Distance != isectAccel.DfGeometry.Distance) {

			msg := fmt.Sprintf(
				"\nDisagreement: hit accel: %t, hit exhaustive: %t\n"+
					"Distance: accel %f, exhaustive %f\n"+
					"Ray: org [%f, %f, %f] dir [%f, %f, %f], mint: %f",
				hitAccel, hitAll,
				isectAccel.DfGeometry.Distance, isectAll.DfGeometry.Distance,
				ray.Origin.X, ray.Origin.Y, ray.Origin.Z,
				ray.Direction.X, ray.Direction.Y, ray.Direction.Z, ray.Mint,
			)

			if hitAll {
				msg += fmt.Sprintf("\nAll hit prim: %d (%s)", isectAll.Primitive.GetID(),
					primitive.GetName(isectAll.Primitive.GetID()))
			}

			if hitAccel {
				msg += fmt.Sprintf("\nAccel hit prim: %d (%s)", isectAccel.Primitive.GetID(),
					primitive.GetName(isectAccel.Primitive.GetID()))
			}

			t.Fatal(msg)
		}
	}
}

func randInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func uniformRandomSphere() geometry.Vector {
	return geometry.NewVector(rand.Float64()-0.5, rand.Float64()-0.5, rand.Float64()-0.5)
}
