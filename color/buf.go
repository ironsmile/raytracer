package color

// Buf implements a weighted buffer of RGB colors. It supports weighted operations for
// updating particular index, reseting the buffer and getting values from the buffer.
type Buf struct {
	pix     []float32
	samples []uint16
}

// Update adds this color into the buffer's ind position.
func (b *Buf) Update(ind int, clr Color) {
	b.UpdateRGB(ind, float32(clr.Red()), float32(clr.Green()), float32(clr.Blue()))
}

// UpdateRGB adds this color into the buffer's ind position.
func (b *Buf) UpdateRGB(ind int, red, green, blue float32) {
	samples := b.samples[ind]
	oldWeight := float32(samples) / float32(samples+1)
	newWeight := 1 - oldWeight

	buffInd := ind * 3
	b.pix[buffInd] = b.pix[buffInd]*oldWeight + red*newWeight
	b.pix[buffInd+1] = b.pix[buffInd+1]*oldWeight + green*newWeight
	b.pix[buffInd+2] = b.pix[buffInd+2]*oldWeight + blue*newWeight

	b.samples[ind]++
}

// Get returns the color at ind in the buffer.
func (b *Buf) Get(ind int) Color {
	buffInd := ind * 3
	return Color{
		red:   float64(b.pix[buffInd]),
		green: float64(b.pix[buffInd+1]),
		blue:  float64(b.pix[buffInd+2]),
	}
}

// Clear erases the sampling counts for the buffer. Update on any index after `Clear` would
// make its value a solid color with weight 1. Note that the buffer is not erased after this
// operation so getting its contents is still possible after Clear.
func (b *Buf) Clear() {
	for i := 0; i < len(b.samples); i++ {
		b.samples[i] = 0
	}
}

// NewBuf returns a buffer with a certain size. The size is the number of
// pixels in the buffer.
func NewBuf(size int) Buf {
	return Buf{
		pix:     make([]float32, size*3),
		samples: make([]uint16, size),
	}
}
