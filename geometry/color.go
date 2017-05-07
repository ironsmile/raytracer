package geometry

import (
	"github.com/ironsmile/raytracer/utils"
)

type Color struct {
	red   float32
	green float32
	blue  float32
}

func (c *Color) clamp() {
	c.red = utils.Clamp32(c.red, 0, 1)
	c.green = utils.Clamp32(c.green, 0, 1)
	c.blue = utils.Clamp32(c.blue, 0, 1)
}

func (c *Color) Red() float32 {
	return c.red
}

func (c *Color) Green() float32 {
	return c.green
}

func (c *Color) Blue() float32 {
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
	s := float32(sclr)
	r := &Color{c.red * s, c.green * s, c.blue * s}
	r.clamp()
	return r
}

func (c *Color) MultiplyScalarIP(sclr float64) *Color {
	s := float32(sclr)
	c.red, c.green, c.blue = c.red*s, c.green*s, c.blue*s
	c.clamp()
	return c
}

func (c *Color) Set(red, green, blue float32) {
	c.red, c.green, c.blue = red, green, blue
	c.clamp()
}

func NewColor(r, g, b float32) *Color {
	return &Color{r, g, b}
}
