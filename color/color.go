package color

import (
	"github.com/ironsmile/raytracer/utils"
)

type Color struct {
	red   float64
	green float64
	blue  float64
}

var (
	// Black is simply the black color
	Black = Color{
		red:   0,
		green: 0,
		blue:  0,
	}
)

func (c *Color) clamp() {
	c.red = utils.Clamp(c.red, 0, 1)
	c.green = utils.Clamp(c.green, 0, 1)
	c.blue = utils.Clamp(c.blue, 0, 1)
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
	r := &Color{c.red + other.red, c.green + other.green, c.blue + other.blue}
	r.clamp()
	return r
}

func (c *Color) PlusIP(other *Color) *Color {
	c.red, c.green, c.blue = c.red+other.red, c.green+other.green, c.blue+other.blue
	c.clamp()
	return c
}

func (c *Color) Multiply(other *Color) *Color {
	return &Color{c.red * other.red, c.green * other.green, c.blue * other.blue}
}

func (c *Color) MultiplyIP(other *Color) *Color {
	c.red, c.green, c.blue = c.red*other.red, c.green*other.green, c.blue*other.blue
	return c
}

func (c *Color) MultiplyScalar(sclr float64) *Color {
	r := &Color{c.red * sclr, c.green * sclr, c.blue * sclr}
	r.clamp()
	return r
}

func (c *Color) MultiplyScalarIP(sclr float64) *Color {
	c.red, c.green, c.blue = c.red*sclr, c.green*sclr, c.blue*sclr
	c.clamp()
	return c
}

func (c *Color) Set(red, green, blue float64) {
	c.red, c.green, c.blue = red, green, blue
	c.clamp()
}

func NewColor(r, g, b float64) *Color {
	return &Color{r, g, b}
}
