package color

import (
	"github.com/ironsmile/raytracer/utils"
)

// Color is a type which represents a RGB color with 3 float values. It also implements
// the standard library's color.Color interface.
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

// Red returns the red color's intensity as a float
func (c *Color) Red() float64 {
	return c.red
}

// Green returns the green color's intensity as a float
func (c *Color) Green() float64 {
	return c.green
}

// Blue returns the blue color's intensity as a float
func (c *Color) Blue() float64 {
	return c.blue
}

// RGBA implements the color.Color interface for this type
func (c *Color) RGBA() (r, g, b, a uint32) {
	return uint32(c.red * 65535), uint32(c.green * 65535), uint32(c.blue * 65535), 65535
}

// Plus "adds" two colors together and returns the result
func (c *Color) Plus(other *Color) *Color {
	r := &Color{c.red + other.red, c.green + other.green, c.blue + other.blue}
	r.clamp()
	return r
}

// PlusIP means "Plus In Place". Its "adds" the other (pass as an argument) color into the
// current (represented by the point receiver) one.
func (c *Color) PlusIP(other *Color) *Color {
	c.red, c.green, c.blue = c.red+other.red, c.green+other.green, c.blue+other.blue
	c.clamp()
	return c
}

// Multiply does a float multiplications between each component for two colors and returns
// the result as a new color.
func (c *Color) Multiply(other *Color) *Color {
	return &Color{c.red * other.red, c.green * other.green, c.blue * other.blue}
}

// MultiplyIP does the same as "Multiply" but the result is stored in the color denoted
// in the point receiver.
func (c *Color) MultiplyIP(other *Color) *Color {
	c.red, c.green, c.blue = c.red*other.red, c.green*other.green, c.blue*other.blue
	return c
}

// MultiplyScalar multiplies every component of a color with a scalar value and returns
// the result as a new color.
func (c *Color) MultiplyScalar(sclr float64) *Color {
	r := &Color{c.red * sclr, c.green * sclr, c.blue * sclr}
	r.clamp()
	return r
}

// MultiplyScalarIP does the same as MultiplyScalar but stores the result in the current
// color.
func (c *Color) MultiplyScalarIP(sclr float64) *Color {
	c.red, c.green, c.blue = c.red*sclr, c.green*sclr, c.blue*sclr
	c.clamp()
	return c
}

// NewColor returns a new color with the supplied in the arguments red, green and blue
// components.
func NewColor(r, g, b float64) *Color {
	return &Color{r, g, b}
}
