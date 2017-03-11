package primitive

import (
	"fmt"

	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/shape"
	"github.com/ironsmile/raytracer/transform"
)

type Rectangle struct {
	BasePrimitive

	objToWorld transform.Transform
	worldToObj transform.Transform
}

func (r *Rectangle) GetType() int {
	return RECTANGLE
}

func NewRectangle(w, h float64) *Rectangle {
	rec := &Rectangle{}
	rec.shape = shape.NewRectangle(w, h)

	rec.objToWorld = *transform.Scale(3, 3, 3).Multiply(transform.RotateY(
		90,
	)).Multiply(transform.Translate(
		geometry.NewVector(0, 0, 0.001),
	))
	rec.worldToObj = *rec.objToWorld.Inverse()

	return rec
}

func (r *Rectangle) Intersect(ray *geometry.Ray, dist float64) (int, float64, *geometry.Vector) {
	objRay := r.worldToObj.Ray(ray)

	res, hitDist, normal := r.shape.Intersect(objRay, dist)

	if res != shape.HIT {
		return res, hitDist, normal
	}

	return res, hitDist, r.objToWorld.Normal(normal)
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle<%s>", r.Name)
}
