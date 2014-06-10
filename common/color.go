package common

type Color Vector

func (c *Color) Red() float64 {
	return c.x
}

func (c *Color) Green() float64 {
	return c.y
}

func (c *Color) Blue() float64 {
	return c.z
}

func (c *Color) RGBA() (r, g, b, a uint32) {
	return uint32(c.x * 65535), uint32(c.y * 65535), uint32(c.z * 65535), 65535
}

func (c *Color) Vector() *Vector {
	return (*Vector)(c)
}

func NewColor(r, g, b float64) *Color {
	col := new(Color)
	col.x = r
	col.y = g
	col.z = b
	return col
}
