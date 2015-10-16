package film

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
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

	// go g.RenderRoutine()

	return nil
}

func (g *GlWindow) RenderRoutine() {

	renderStart := time.Now()

	g.window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	fmt.Printf("gl.GetError: %d\n", gl.GetError())
	fmt.Printf("gl.VERSION = %s\n", gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Printf("gl.RENDERER = %s\n", gl.GoStr(gl.GetString(gl.RENDERER)))
	fmt.Printf("gl.VENDOR = %s\n", gl.GoStr(gl.GetString(gl.VENDOR)))

	// gl.MatrixMode(gl.PROJECTION)
	// gl.LoadIdentity()
	// gl.Ortho(0, float64(g.width), float64(g.height), 0, 0, 1)
	gl.Disable(gl.DEPTH_TEST)
	// gl.MatrixMode(gl.MODELVIEW)
	// gl.LoadIdentity()
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	var texture uint32
	gl.GenTextures(1, &texture)

	// gl.PushAttrib(gl.ENABLE_BIT)
	gl.Enable(gl.TEXTURE_2D)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	// defer func() {
	// 	g.window.MakeContextCurrent()
	// 	// texture.Unbind(gl.TEXTURE_2D)
	// 	// gl.PopAttrib()
	// 	gl.Disable(gl.TEXTURE_2D)
	// }()

	displayTexture := func() {
		// textureTime := time.Now()

		g.pixBufferLock.RLock()
		defer g.pixBufferLock.RUnlock()

		g.window.MakeContextCurrent()

		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(g.width), int32(g.height),
			0, gl.RGB, gl.FLOAT, gl.Ptr(g.pixBuffer))

		//!TODO: Somehow render the preparte texture.
		// gl.Begin(gl.POLYGON)

		// gl.TexCoord2xOES(0, 0)
		// gl.Vertex3xOES(0, 0)

		// gl.TexCoord2xOES(1, 0)
		// gl.Vertex3xOES(int32(g.width), 0)

		// gl.TexCoord2xOES(1, 1)
		// gl.Vertex3xOES(int32(g.width), int32(g.height))

		// gl.TexCoord2xOES(0, 1)
		// gl.Vertex3xOES(0, int32(g.height))

		// gl.End()

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

	// This select makes Wait reentrant. After the first Wait all others will just return
	// at once. In other words, all Waits after the first one are a noop.
	select {
	case <-g.renderFinishChan:
		return
	default:
		fmt.Println("Sending rendering finish")
		g.renderFinishChan <- true

		fmt.Println("Receiving finished ack")
		_ = <-g.renderFinishChan

		fmt.Println("Closing renderFinishChan")
		close(g.renderFinishChan)

		fmt.Println("Closing refreshScreenChan")
		close(g.refreshScreenChan)
	}
}

func (g *GlWindow) StartFrame() {
}

func (g *GlWindow) DoneFrame() {
	select {
	case <-g.renderFinishChan:
		return
	default:
		g.RefreshScreen()
	}
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
