package scene

import (
	"github.com/ironsmile/raytracer/geometry"
)

type Object struct {
	BasePrimitive

	Triangles []Primitive
}

func (o *Object) GetType() int {
	return OBJECT
}

func (o *Object) Intersect(ray *geometry.Ray, dist float64) (int, float64) {
	prim, distance := IntersectPrimitives(o.Triangles, ray)
	if prim == nil {
		return MISS, distance
	}
	return HIT, distance
}
