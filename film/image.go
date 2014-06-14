package film

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Image struct {
	width  int
	height int

	img *image.NRGBA

	filename string
}

func (i *Image) Wait() {

}

func (i *Image) Init(width int, height int) error {
	i.width = width
	i.height = height

	i.img = image.NewNRGBA(image.Rect(0, 0, width, height))

	return nil
}

func (i *Image) Width() int {
	return i.width
}

func (i *Image) Height() int {
	return i.height
}

func (i *Image) Done() {
	out, err := os.Create(i.filename)

	if err != nil {
		fmt.Errorf("%s\n", err.Error())
	}

	defer func() {
		out.Close()
	}()

	err = png.Encode(out, i.img)

	if err != nil {
		fmt.Errorf("%s\n", err.Error())
	}
}

func (i *Image) Set(x, y int, clr color.Color) error {
	i.img.Set(x, y, clr)
	return nil
}

func NewImage(filname string) *Image {
	img := new(Image)
	img.filename = filname
	return img
}
