package accel

import (
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/primitive"
)

// Voxel structure records which primitives overlap its extent using a vector
type Voxel struct {
	primitives      []primitive.Primitive
	allCanIntersect bool
}

// NewVoxel returns a new initialized voxel
func NewVoxel() *Voxel {
	return &Voxel{allCanIntersect: true}
}

// Intersect checks whether a ray itnersects any of the voxel's primitives
func (v *Voxel) Intersect(ray geometry.Ray, in *primitive.Intersection) bool {
	return primitive.IntersectMultiple(v.primitives, ray, in)
}

// IntersectP checks whether a ray itnersects any of the voxel's primitives. It does not
// return intersection data so it is faster than the Intersect method.
func (v *Voxel) IntersectP(ray geometry.Ray) bool {
	return primitive.IntersectPMultiple(v.primitives, ray)
}

// Add inserts a primitive in this voxel
func (v *Voxel) Add(p primitive.Primitive) {
	if p.CanIntersect() {
		v.primitives = append(v.primitives, p)
	} else {
		gridAccel := NewGrid(p.Refine())
		v.primitives = append(v.primitives, gridAccel)
	}
}
