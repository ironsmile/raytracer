package film

import "image/color"

type Film interface {
	Width() int
	Height() int

	Init(width int, height int) error
	Set(x int, y int, clr color.Color) error
	Ping()
	Done()
	Wait()
}
