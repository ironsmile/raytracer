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

	pixBuffer := make([]float32, g.width*g.height*3)

	g.window.MakeContextCurrent()

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Viewport(0, 0, g.width, g.height)
	gl.Ortho(0, float64(g.width), float64(g.height), 0, 0, 1)
	gl.Disable(gl.DEPTH_TEST)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Clear(gl.COLOR_BUFFER_BIT)

	defer func() {
		g.window.MakeContextCurrent()

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

		// texture.Unbind(gl.TEXTURE_2D)
		// gl.PopAttrib()

		// gl.PushAttrib(gl.ENABLE_BIT)
		// gl.Enable(gl.TEXTURE_2D)
		// texture.Bind(gl.TEXTURE_2D)

		fmt.Println("O? Why TexImage2D? Why did you do this?")
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, g.width, g.height, 0, gl.RGB, gl.FLOAT,
			pixBuffer)
		// fmt.Println("MipMap, was it you then?")
		// gl.GenerateMipmap(gl.TEXTURE_2D)

		fmt.Println("Ohhh here it goes!")

		texture.Unbind(gl.TEXTURE_2D)
		gl.PopAttrib()

		fmt.Println("Beware. Will activate texture")

		// gl.ActiveTexture(gl.TEXTURE0)

		fmt.Println("Will enable texture_2d?")
		gl.Enable(gl.TEXTURE_2D)

		fmt.Println("Will bind texture?")
		texture.Bind(gl.TEXTURE_2D)

		fmt.Println("Textured polygon!")
		gl.Begin(gl.POLYGON)

		gl.TexCoord2f(0, 1)
		gl.Vertex2i(0, 0)

		gl.TexCoord2f(1, 1)
		gl.Vertex2i(g.width, 0)

		gl.TexCoord2f(1, 0)
		gl.Vertex2i(g.width, g.height)

		gl.TexCoord2f(0, 0)
		gl.Vertex2i(0, g.height)

		gl.End()

		fmt.Println("Polygon")

		gl.Disable(gl.TEXTURE_2D)

		fmt.Println("GL rendering goroutine exited.")
		fmt.Printf("GL Points drawing: %s\n", pointsTime)
		fmt.Printf("GL SwapBuffers: %s\n", swapBuffersTime)

		g.window.SwapBuffers()
	}()

	renderPixel := func(pixel *PixelInfo) {
		pixBuffer[g.width*pixel.GlMatrixY*3+pixel.GlMatrixX*3] = pixel.RedFloat
		pixBuffer[g.width*pixel.GlMatrixY*3+pixel.GlMatrixX*3+1] = pixel.GreenFloat
		pixBuffer[g.width*pixel.GlMatrixY*3+pixel.GlMatrixX*3+2] = pixel.BlueFloat
		// gl.Color3f(pixel.RedFloat, pixel.GreenFloat, pixel.BlueFloat)
		// gl.Vertex2i(pixel.GlMatrixX, pixel.GlMatrixY)
	}

	fmt.Printf("GL Init time: %s\n", time.Since(renderStart))

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
