package geometry

type Ray struct {
	Origin    *Point
	Direction *Vector

	Mint float64
	Maxt float64

	Time float64

	Depth int

	Debug bool
}

func (r *Ray) AtTime(time float64) *Point {
	return r.Origin.PlusVector(r.Direction.MultiplyScalar(time))
}

func NewRay(origin Point, dir Vector) *Ray {
	return &Ray{Origin: &origin, Direction: &dir}
}

func NewRayFull(origin Point, dir Vector, start, end float64) *Ray {
	return &Ray{Origin: &origin, Direction: &dir, Mint: start, Maxt: end}
}
