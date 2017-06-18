package film

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	rColor "github.com/ironsmile/raytracer/color"
)

type Image struct {
	width  int
	height int

	buf rColor.Buf

	filename string
}

func (i *Image) Wait() {

}

func (i *Image) Init(width int, height int) error {
	i.width = width
	i.height = height
	i.buf = rColor.NewBuf(width * height)
	return nil
}

func (i *Image) Width() int {
	return i.width
}

func (i *Image) Height() int {
	return i.height
}

func (i *Image) DoneFrame() {
	out, err := os.Create(i.filename)

	if err != nil {
		fmt.Errorf("%s\n", err.Error())
	}

	defer func() {
		out.Close()
	}()

	img := image.NewNRGBA(image.Rect(0, 0, i.width, i.height))
	buffSize := i.width * i.height
	for ind := 0; ind < buffSize; ind++ {
		y := ind / i.width
		x := ind % i.width
		bufClr := i.buf.Get(ind)
		img.Set(x, y, color.NRGBA64{
			R: uint16(bufClr.Red() * float64(math.MaxUint16)),
			G: uint16(bufClr.Green() * float64(math.MaxUint16)),
			B: uint16(bufClr.Blue() * float64(math.MaxUint16)),
			A: math.MaxUint16,
		})
	}

	err = png.Encode(out, img)

	if err != nil {
		fmt.Printf("encoding image failed: %s\n", err)
	} else {
		fmt.Printf("Image saved to %s\n", i.filename)
	}
}

func (i *Image) StartFrame() {}

func (i *Image) Set(x, y int, clr color.Color) error {
	ind := i.width*y + x
	ri, gi, bi, _ := clr.RGBA()
	i.buf.UpdateRGB(ind, float32(ri)/65535.0, float32(gi)/65535.0, float32(bi)/65535.0)
	return nil
}

func NewImage(filname string) *Image {
	img := new(Image)
	img.filename = filname
	return img
}
