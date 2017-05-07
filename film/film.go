package film

import "github.com/ironsmile/raytracer/geometry"

type Film interface {
	Width() int
	Height() int

	Init(width int, height int) error
	Set(x int, y int, clr *geometry.Color) error

	StartFrame()
	DoneFrame()
	Wait()
}
