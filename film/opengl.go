package film

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/go-gl-legacy/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type GlWindow struct {
	width  int
	height int

	refreshScreenChan chan bool
	renderFinishChan  chan bool

	window *glfw.Window

	pixBufferLock sync.RWMutex
	pixBuffer     []float32
}

func (g *GlWindow) Init(width int, height int) error {
	g.width = width
	g.height = height

	g.pixBuffer = make([]float32, g.width*g.height*3)
	g.refreshScreenChan = make(chan bool)
	g.renderFinishChan = make(chan bool)

	go g.renderRoutine()

	return nil
}

func (g *GlWindow) renderRoutine() {

	renderStart := time.Now()

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

	defer func() {
		g.window.MakeContextCurrent()
		texture.Unbind(gl.TEXTURE_2D)
		gl.PopAttrib()
		gl.Disable(gl.TEXTURE_2D)
	}()

	fmt.Printf("gl.RENDERER = %s\n", gl.GetString(gl.RENDERER))
	fmt.Printf("gl.VERSION = %s\n", gl.GetString(gl.VERSION))
	fmt.Printf("gl.VENDOR = %s\n", gl.GetString(gl.VENDOR))

	displayTexture := func() {
		// textureTime := time.Now()

		g.pixBufferLock.RLock()
		defer g.pixBufferLock.RUnlock()

		g.window.MakeContextCurrent()

		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, g.width, g.height, 0, gl.RGB, gl.FLOAT,
			g.pixBuffer)

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

		g.window.SwapBuffers()

		// fmt.Printf("GL Textured Polygon: %s\n", time.Since(textureTime))
	}

	defer func() {
		fmt.Println("GL rendering goroutine exited.")
	}()

	fmt.Printf("GL Init time: %s\n", time.Since(renderStart))

	for {

		select {
		case _ = <-g.refreshScreenChan:
			g.refreshScreenChan <- true
			displayTexture()

		case _ = <-g.renderFinishChan:
			g.renderFinishChan <- true
			return
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
}

func (g *GlWindow) StartFrame() {
}

func (g *GlWindow) DoneFrame() {
	g.RefreshScreen()
}

func (g *GlWindow) Set(x int, y int, clr color.Color) error {
	g.pixBufferLock.Lock()
	defer g.pixBufferLock.Unlock()

	ri, gi, bi, _ := clr.RGBA()

	ind := g.width*y*3 + x*3
	g.pixBuffer[ind] = float32(ri) / 65535.0
	g.pixBuffer[ind+1] = float32(gi) / 65535.0
	g.pixBuffer[ind+2] = float32(bi) / 65535.0

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

func NewGlWIndow(window *glfw.Window) *GlWindow {
	gwWin := &GlWindow{window: window}
	return gwWin
}
