package color

type pixel struct {
	samples int
	c, v    Color
}

// Buf implements a weighted buffer of RGB colors. It supports weighted operations for
// updating particular index, reseting the buffer and getting values from the buffer.
type Buf []pixel

// Update adds this color into the buffer's ind position.
func (b Buf) Update(ind int, clr Color) {
	p := &[]pixel(b)[ind]
	p.samples++
	if p.samples == 1 {
		p.c = clr
		return
	}
	t := p.c
	p.c = *p.c.Plus(clr.Minus(&p.c).MultiplyScalar(1 / float64(p.samples)))
	p.v = *p.v.Plus(clr.Minus(&t).Multiply(clr.Minus(&p.c)))
}

// UpdateRGB adds this color into the buffer's ind position.
func (b Buf) UpdateRGB(ind int, red, green, blue float32) {
	b.Update(ind, Color{
		red:   float64(red),
		green: float64(green),
		blue:  float64(blue),
	})
}

// Get returns the color at ind in the buffer.
func (b Buf) Get(ind int) Color {
	pix := []pixel(b)[ind]
	return *pix.c.MultiplyScalar(1 / float64(pix.samples)).Clamp()
}

// Clear erases the sampling counts for the buffer. Update on any index after `Clear` would
// make its value a solid color with weight 1. Note that the buffer is not erased after this
// operation so getting its contents is still possible after Clear.
func (b Buf) Clear() {
	bs := []pixel(b)
	for i := 0; i < len(bs); i++ {
		bs[i].samples = 0
		bs[i].v = Black
	}
}

// NewBuf returns a buffer with a certain size. The size is the number of
// pixels in the buffer.
func NewBuf(size int) Buf {
	return make(Buf, size)
}
