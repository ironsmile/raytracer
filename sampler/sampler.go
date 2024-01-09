package sampler

import "image/color"

// Output supports setting colours to certain 2D pixels.
type Output interface {
	Set(x int, y int, clr color.Color) error

	DoneFrame()
	StartFrame()
}
