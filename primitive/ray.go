package primitive

import (
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/mat"
	"github.com/ironsmile/raytracer/transform"
)

// Ray represents a ray in the world. It has a particular "width" along its direction.
//
// TODO: explore implementation based on the [geometry.Ray.Intersect] function.
type Ray struct {
	BasePrimitive

	ray geometry.Ray
}

// NewRay returns a new Ray [Primitive] which could be seen in the world.
func NewRay(ray geometry.Ray) *Ray {
	r := &Ray{
		ray: ray,
	}
	r.SetTransform(transform.Identity())
	r.id = GetNewID()
	return r
}

// CanIntersect returns false as the ray is composed of few other primitives.
func (r *Ray) CanIntersect() bool {
	return false
}

// Refine returns the primitives from which this ray is composed of.
func (r *Ray) Refine() []Primitive {
	rayMat := mat.NewMaterial()
	rayMat.Color = geometry.NewColor(1, 1, 0)
	rayMat.Diff = 1

	origin := NewSphere(0.04)
	origin.SetTransform(
		transform.Translate(r.ray.Origin),
	)
	origin.Shape().SetMaterial(*rayMat)

	rayEnd := r.ray.Origin.Plus(
		r.ray.Direction.Normalize().MultiplyScalar(r.ray.Maxt),
	)

	endMat := mat.NewMaterial()
	endMat.Color = geometry.NewColor(0, 1, 1)
	endMat.Diff = 1
	end := NewSphere(0.04)
	end.SetTransform(
		transform.Translate(rayEnd),
	)
	end.Shape().SetMaterial(*endMat)

	line := NewCylinder(0.03, r.ray.Origin, rayEnd)
	line.Shape().SetMaterial(*rayMat)

	return []Primitive{
		origin,
		line,
		end,
	}
}
