package geometry

type Ray struct {
	Origin    *Point
	Direction *Vector
	Debug     bool
}

func NewRay(origin Point, dir Vector) *Ray {
	return &Ray{Origin: &origin, Direction: &dir}
}
