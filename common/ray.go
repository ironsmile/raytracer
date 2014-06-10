package common

type Ray struct {
	Origin    *Vector
	Direction *Vector
	Debug     bool
}

func NewRay(origin, dir Vector) *Ray {
	ray := new(Ray)
	ray.Origin = &origin
	ray.Direction = &dir
	return ray
}
