package geometry

type Color struct {
	red   float64
	green float64
	blue  float64
}

func (c *Color) Red() float64 {
	return c.red
}

func (c *Color) Green() float64 {
	return c.green
}

func (c *Color) Blue() float64 {
	return c.blue
}

func (c *Color) RGBA() (r, g, b, a uint32) {
	return uint32(c.red * 65535), uint32(c.green * 65535), uint32(c.blue * 65535), 65535
}

func (c *Color) Plus(other *Color) *Color {
	return &Color{c.red + other.red, c.green + other.green, c.blue + other.blue}
}

func (c *Color) Multiply(other *Color) *Color {
	return &Color{c.red * other.red, c.green * other.green, c.blue * other.blue}
}

func (c *Color) MultiplyScalar(sclr float64) *Color {
	return &Color{c.red * sclr, c.green * sclr, c.blue * sclr}
}

func NewColor(r, g, b float64) *Color {
	return &Color{r, g, b}
}
