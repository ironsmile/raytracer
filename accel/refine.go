package accel

import (
	"github.com/ironsmile/raytracer/primitive"
)

// FullyRefinePrimitives gets a slice of primitives which may be intersectable or not. It
// then refines all primitives in the list recursively untill they are all intersectable
// and returns a list with only intersectable primitives.
func FullyRefinePrimitives(p []primitive.Primitive) []primitive.Primitive {
	var refined []primitive.Primitive

	// Fully refine the primitives and add them in the bvh.primitives
	todo := NewPriorityQueue(p)

	for todo.Len() > 0 {
		prim := todo.PopPrimitive()
		if prim.CanIntersect() {
			refined = append(refined, prim)
			continue
		}

		for _, deg := range prim.Refine() {
			if deg.CanIntersect() {
				refined = append(refined, deg)
			} else {
				todo.PushPrimitive(deg)
			}
		}
	}

	return refined
}
