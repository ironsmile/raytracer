package film

import (
	"fmt"
	"image/color"
	"time"

	"github.com/go-gl/gl"
	"github.com/go-gl/glfw3"
)

type PixelInfo struct {
	RedFloat   float32
	GreenFloat float32
	BlueFloat  float32

	GlMatrixX int
	GlMatrixY int
}

type GlWindow struct {
	width  int
	height int

	renderChan        chan *PixelInfo
	refreshScreenChan chan bool
	renderFinishChan  chan bool
	frameFinishChan   chan bool
	renderTasksChan   chan chan *PixelInfo

	window *glfw3.Window
}

func (g *GlWindow) Init(width int, height int) error {
	g.width = width
	g.height = height

	g.refreshScreenChan = make(chan bool)
	g.renderFinishChan = make(chan bool)
	g.frameFinishChan = make(chan bool)
	g.renderTasksChan = make(chan chan *PixelInfo)

	g.window.MakeContextCurrent()
	g.window.SwapBuffers()

	go g.renderRoutine()

	return nil
}

func (g *GlWindow) renderRoutine() {

	var pointsTime time.Duration
	var refreshTime time.Duration
	renderStart := time.Now()

	timed := func(fnc func()) time.Duration {
		s := time.Now()
		fnc()
		return time.Since(s)
	}

	pixBuffer := make([]float32, g.width*g.height*3)

	g.window.MakeContextCurrent()

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(g.width), float64(g.height), 0, 0, 1)
	gl.Disable(gl.DEPTH_TEST)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	texture := gl.GenTexture()

	gl.PushAttrib(gl.ENABLE_BIT)
	gl.Enable(gl.TEXTURE_2D)
	texture.Bind(gl.TEXTURE_2D)

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	texture.Unbind(gl.TEXTURE_2D)
	gl.PopAttrib()
	gl.Disable(gl.TEXTURE_2D)

	displayTexture := func() {
		textureTime := time.Now()

		g.window.MakeContextCurrent()

		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.PushAttrib(gl.ENABLE_BIT)
		gl.Enable(gl.TEXTURE_2D)
		texture.Bind(gl.TEXTURE_2D)

		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, g.width, g.height, 0, gl.RGB, gl.FLOAT,
			pixBuffer)

		gl.Begin(gl.POLYGON)

		gl.TexCoord2f(0, 0)
		gl.Vertex2i(0, 0)

		gl.TexCoord2f(1, 0)
		gl.Vertex2i(g.width, 0)

		gl.TexCoord2f(1, 1)
		gl.Vertex2i(g.width, g.height)

		gl.TexCoord2f(0, 1)
		gl.Vertex2i(0, g.height)

		gl.End()

		texture.Unbind(gl.TEXTURE_2D)
		gl.PopAttrib()
		gl.Disable(gl.TEXTURE_2D)

		g.window.SwapBuffers()

		fmt.Printf("GL Textured Polygon: %s\n", time.Since(textureTime))
		fmt.Printf("GL Points drawing: %s\n", pointsTime)
	}

	defer func() {
		fmt.Println("GL rendering goroutine exited.")
		fmt.Printf("GL screen refreshes: %s\n", refreshTime)

	}()

	addPixelToBuffer := func(pixel *PixelInfo) {
		ind := g.width*pixel.GlMatrixY*3 + pixel.GlMatrixX*3
		pixBuffer[ind] = pixel.RedFloat
		pixBuffer[ind+1] = pixel.GreenFloat
		pixBuffer[ind+2] = pixel.BlueFloat
	}

	fmt.Printf("GL Init time: %s\n", time.Since(renderStart))

	renderChan := <-g.renderTasksChan
	pointsTime = 0

	for {

		select {
		// case pInfo := <-renderChan:

		// 	pointsTime += timed(func() {
		// 		if pInfo != nil {
		// 			addPixelToBuffer(pInfo)
		// 		}
		// 	})

		case _ = <-g.refreshScreenChan:

			refreshTime += timed(func() {
				displayTexture()
			})
			g.refreshScreenChan <- true

		case _ = <-g.renderFinishChan:
			g.renderFinishChan <- true
			return

		case _ = <-g.frameFinishChan:

			for pInfo := range renderChan {
				pointsTime += timed(func() {
					addPixelToBuffer(pInfo)
				})
			}
			g.frameFinishChan <- true
			displayTexture()

		case renderChan = <-g.renderTasksChan:
			pointsTime = 0
		}

	}
}

func (g *GlWindow) Wait() {

	fmt.Println("Sending rendering finish")
	g.renderFinishChan <- true

	fmt.Println("Receiving finished ack")
	_ = <-g.renderFinishChan

	fmt.Println("Closing refreshScreenChan")
	close(g.refreshScreenChan)

	fmt.Println("Closing renderFinishChan")
	close(g.renderFinishChan)

	fmt.Println("Closing frameFinishChan")
	close(g.frameFinishChan)

	fmt.Println("Closing renderTasksChan")
	close(g.renderTasksChan)
}

func (g *GlWindow) StartFrame() {
	chanBuffer := g.width * g.height
	if chanBuffer > 1e7 {
		chanBuffer = 1e7
	}

	g.renderChan = make(chan *PixelInfo, chanBuffer)
	g.renderTasksChan <- g.renderChan
}

func (g *GlWindow) DoneFrame() {
	g.frameFinishChan <- true
	close(g.renderChan)
	_ = <-g.frameFinishChan
}

func (g *GlWindow) Set(x int, y int, clr color.Color) error {

	ri, gi, bi, _ := clr.RGBA()

	pInfo := &PixelInfo{
		float32(ri) / 65535.0,
		float32(gi) / 65535.0,
		float32(bi) / 65535.0,
		x,
		y}

	g.renderChan <- pInfo

	return nil
}

func (g *GlWindow) RefreshScreen() {
	g.refreshScreenChan <- true
	_ = <-g.refreshScreenChan
}

func (g *GlWindow) Width() int {
	return g.width
}

func (g *GlWindow) Height() int {
	return g.height
}

func NewGlWIndow(window *glfw3.Window) *GlWindow {
	gwWin := &GlWindow{window: window}
	return gwWin
}
