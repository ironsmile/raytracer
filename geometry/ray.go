package geometry

type Ray struct {
	Origin    Vector
	Direction Vector

	Debug bool
}

func (r *Ray) BackToDefaults() {
	r.Debug = false
}
