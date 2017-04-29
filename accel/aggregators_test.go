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
			accel: NewBVH(prims, 3),
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
	var bb *bbox.BBox

	for _, pr := range prims {
		bb = bbox.Union(bb, pr.GetWorldBBox())
	}

	bb.Expand(bb.Max.ByAxis(bb.MaximumExtend()) -
		bb.Min.ByAxis(bb.MaximumExtend()))

	var lastHit geometry.Vector

	for i := 0; i < 2000; i++ {
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

		var isectAll, isectAccel primitive.Intersection
		hitAll := primitive.IntersectMultiple(prims, ray, &isectAll)
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
				msg += fmt.Sprintf("\nAll hit prim: %d", isectAll.Primitive.GetID())
			}

			if hitAccel {
				msg += fmt.Sprintf("\nAccel hit prim: %d", isectAccel.Primitive.GetID())
			}

			t.Fatal(msg)
		}

		if hitAll {
			lastHit = ray.At(isectAll.DfGeometry.Distance)
		}
	}
}

func randInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func uniformRandomSphere() geometry.Vector {
	return geometry.NewVector(rand.Float64()-0.5, rand.Float64()-0.5, rand.Float64()-0.5)
}
