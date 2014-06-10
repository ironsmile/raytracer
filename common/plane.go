package common

type Plane struct {
	N Vector
	D float64
}

func NewPlane(normal Vector, d float64) *Plane {
	plane := new(Plane)
	plane.N = normal
	plane.D = d
	return plane
}
