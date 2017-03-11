package shape

import "github.com/ironsmile/raytracer/geometry"

type Rectangle struct {
	width  float64
	height float64
}

func NewRectangle(w, h float64) *Rectangle {
	if w < 0 || w > 1 || h < 0 || h > 1 {
		panic("Recatangle width and height must be in the [0-1] region")
	}
	return &Rectangle{w, h}
}

func (r *Rectangle) Intersect(ray *geometry.Ray, dist float64) (int, float64, *geometry.Vector) {
	normal := geometry.NewVector(0, 0, -1)

	d := geometry.NewPoint(0, 0, 0).MinusVectorIP(ray.Origin.Vector()).Vector().Product(normal)
	d /= ray.Direction.Product(normal)

	if d <= 0 {
		return MISS, dist, nil
	}

	hitPoint := ray.Origin.PlusVector(ray.Direction.MultiplyScalar(d))

	if hitPoint.X >= 0 && hitPoint.X <= r.width && hitPoint.Y >= 0 && hitPoint.Y <= r.height {
		return HIT, d, normal.Copy()
	}

	return MISS, dist, nil
}
