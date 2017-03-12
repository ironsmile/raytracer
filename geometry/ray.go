package geometry

type Ray struct {
	Origin    Point
	Direction Vector

	Debug bool
}

func (r *Ray) BackToDefaults() {
	r.Debug = false
}
