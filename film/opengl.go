package film

import (
	"fmt"
	"image/color"
	"sync"

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

	glProgram uint32 // Holds the OpenGL program
	glVao     uint32
	glVbo     uint32
	glEbo     uint32
	glTexture uint32

	glInited bool
}

func (g *GlWindow) Init(width int, height int) error {
	g.width = width
	g.height = height

	g.pixBuffer = make([]float32, g.width*g.height*3)
	g.refreshScreenChan = make(chan bool)
	g.renderFinishChan = make(chan bool)

	return g.initOpenGL()
}

func (g *GlWindow) initOpenGL() error {
	if g.glInited {
		panic("Calling init on already initialized OpenGL")
	}

	g.window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return err
	}

	fmt.Printf("gl.VERSION = %s\n", gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Printf("gl.RENDERER = %s\n", gl.GoStr(gl.GetString(gl.RENDERER)))
	fmt.Printf("gl.VENDOR = %s\n", gl.GoStr(gl.GetString(gl.VENDOR)))

	// General configuration
	gl.Viewport(0, 0, int32(g.width), int32(g.height))
	gl.Disable(gl.DEPTH_TEST)

	// Shaders and program compilation
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		return err
	}

	g.glProgram = program

	// Texture creation
	gl.GenTextures(1, &g.glTexture)
	gl.BindTexture(gl.TEXTURE_2D, g.glTexture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	// Vertex buffer, array
	gl.GenVertexArrays(1, &g.glVao)
	gl.BindVertexArray(g.glVao)
	{

		verices := []float32{
			// [0, 3] - Positions, [4, 5] - Texture cords
			-1, 1, 0, 0, 0, // top left [0]
			1, 1, 0, 1, 0, // top right [1]
			-1, -1, 0, 0, 1, // bottom left [2]
			1, -1, 0, 1, 1, // bottom right [3]
		}

		indeces := []uint32{
			0, 1, 3, // First Triangle
			3, 0, 2, // Second Triangle
		}

		gl.GenBuffers(1, &g.glVbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, g.glVbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(verices)*4, gl.Ptr(verices), gl.STATIC_DRAW)

		gl.GenBuffers(1, &g.glEbo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, g.glEbo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indeces)*4, gl.Ptr(indeces), gl.STATIC_DRAW)

		// Vertex Attribute
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
		gl.EnableVertexAttribArray(0)

		// Texture Position Attribute
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
		gl.EnableVertexAttribArray(1)
	}
	gl.BindVertexArray(0)

	g.glInited = true
	return nil
}

func (g *GlWindow) Render() {
	if !g.glInited {
		panic("OpenGL not initialized!")
	}

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.UseProgram(g.glProgram)
	gl.BindVertexArray(g.glVao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, g.glTexture)

	g.bufferToTexture()
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (g *GlWindow) bufferToTexture() {
	g.pixBufferLock.RLock()
	defer g.pixBufferLock.RUnlock()
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(g.width), int32(g.height),
		0, gl.RGB, gl.FLOAT, gl.Ptr(g.pixBuffer))
}

func (g *GlWindow) Wait() {
	gl.DeleteVertexArrays(1, &g.glVao)
	gl.DeleteBuffers(1, &g.glEbo)
	gl.DeleteBuffers(1, &g.glVbo)
}

func (g *GlWindow) StartFrame() {
}

func (g *GlWindow) DoneFrame() {
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
