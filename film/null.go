package film

import "github.com/ironsmile/raytracer/geometry"

type NullFilm struct {
	width  int
	height int
}

func (n *NullFilm) Set(x int, y int, clr *geometry.Color) error {
	return nil
}

func (n *NullFilm) Width() int {
	return n.width
}

func (n *NullFilm) Height() int {
	return n.height
}

func (n *NullFilm) Init(width int, height int) error {
	n.width = width
	n.height = height
	return nil
}

func (n *NullFilm) DoneFrame() {

}

func (n *NullFilm) StartFrame() {

}

func (n *NullFilm) Wait() {

}

func NewNullFilm() *NullFilm {
	nf := new(NullFilm)
	return nf
}
