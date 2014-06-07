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

	renderChan      chan *PixelInfo
	swapBuffersChan chan bool
	renderStopChan  chan bool

	window *glfw3.Window
}

func (g *GlWindow) Init(width int, height int) error {
	g.width = width
	g.height = height

	hasWindow := glfw3.Init()

	if !hasWindow {
		return fmt.Errorf("Initializing glfw3 failed")
	}

	window, err := glfw3.CreateWindow(g.width, g.height, "Raytracer", nil, nil)

	if err != nil {
		return err
	}

	g.window = window

	g.window.SetCloseCallback(func(w *glfw3.Window) {
		g.window.SetShouldClose(true)
	})

	g.window.SetKeyCallback(func(w *glfw3.Window, key glfw3.Key, scancode int,
		action glfw3.Action, mods glfw3.ModifierKey) {
		if key != glfw3.KeyEscape {
			return
		}
		g.window.SetShouldClose(true)
	})

	chanBuffer := g.width * g.height
	if chanBuffer > 1e7 {
		chanBuffer = 1e7
	}

	g.renderChan = make(chan *PixelInfo, chanBuffer)
	g.swapBuffersChan = make(chan bool)
	g.renderStopChan = make(chan bool)

	g.window.MakeContextCurrent()
	g.window.SwapBuffers()

	go g.renderRoutine()

	return nil
}

func (g *GlWindow) renderRoutine() {

	var pointsTime time.Duration
	var swapBuffersTime time.Duration
	renderStart := time.Now()

	timed := func(fnc func()) time.Duration {
		s := time.Now()
		fnc()
		return time.Since(s)
	}

	defer func() {
		gl.End()
		fmt.Println("GL rendering goroutine exited.")
		fmt.Printf("GL Points drawing: %s\n", pointsTime)
		fmt.Printf("GL SwapBuffers: %s\n", swapBuffersTime)
		g.window.MakeContextCurrent()
		g.window.SwapBuffers()
	}()

	g.window.MakeContextCurrent()

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(g.width), float64(g.height), 0, 0, 1)
	gl.Disable(gl.DEPTH_TEST)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(0.375, 0.375, 0)

	// gl.DrawPixels(width, height, format, typ, pixels)
	// b := gl.GenBuffer()
	// b.Bind(gl.ARRAY_BUFFER)

	renderPixel := func(pixel *PixelInfo) {
		gl.Color3f(pixel.RedFloat, pixel.GreenFloat, pixel.BlueFloat)
		gl.Vertex2i(pixel.GlMatrixX, pixel.GlMatrixY)
	}

	fmt.Printf("GL Init time: %s\n", time.Since(renderStart))

	gl.Begin(gl.POINTS)
	for {

		select {
		case pInfo := <-g.renderChan:
			g.window.MakeContextCurrent()

			pointsTime += timed(func() {
				renderPixel(pInfo)
			})

		case _ = <-g.swapBuffersChan:
			swapBuffersTime += timed(func() {
				g.window.SwapBuffers()
			})
			g.swapBuffersChan <- true
		case _ = <-g.renderStopChan:
			g.window.MakeContextCurrent()

			for pInfo := range g.renderChan {
				pointsTime += timed(func() {
					renderPixel(pInfo)
				})
			}
			g.renderStopChan <- true

			return
		}

	}
}

func (g *GlWindow) closeWindow() {
	close(g.swapBuffersChan)
	close(g.renderStopChan)
	g.window.Destroy()
	glfw3.Terminate()
}

func (g *GlWindow) Done() {
	g.renderStopChan <- true
	close(g.renderChan)
	_ = <-g.renderStopChan
}

func (g *GlWindow) Wait() {
	for !g.window.ShouldClose() {
		glfw3.WaitEvents()
	}

	g.closeWindow()
}

func (g *GlWindow) Set(x int, y int, clr color.Color) error {

	ri, gi, bi, _ := clr.RGBA()

	pInfo := new(PixelInfo)

	pInfo.RedFloat = float32(ri) / 65535.0
	pInfo.GreenFloat = float32(gi) / 65535.0
	pInfo.BlueFloat = float32(bi) / 65535.0

	pInfo.GlMatrixX = x
	pInfo.GlMatrixY = y

	g.renderChan <- pInfo

	return nil
}

func (g *GlWindow) Ping() {
	g.swapBuffersChan <- true
	_ = <-g.swapBuffersChan
}

func (g *GlWindow) Width() int {
	return g.width
}

func (g *GlWindow) Height() int {
	return g.height
}

func NewGlWIndow() *GlWindow {
	win := new(GlWindow)
	return win
}
