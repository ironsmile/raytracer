package film

import (
	"fmt"
	"image/color"
	"os"
	"runtime/trace"
	"sync"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ironsmile/raytracer/camera"
	"github.com/ironsmile/raytracer/engine"
	"github.com/ironsmile/raytracer/sampler"
	"github.com/ironsmile/raytracer/scene"
)

// GlWinArgs are passed to the GlWindow.Run() and control some aspects on how
// the application would be ran.
type GlWinArgs struct {
	Fullscreen  bool
	VSync       bool
	Width       int
	Height      int
	Interactive bool
	ShowBBoxes  bool
	FPSCap      uint
	ShowFPS     bool
	SceneName   string
}

// GlWindow is a `film` which renders the scene in an OpenGL window using GLFW3.
// The scene is renderd in a texture which is applied on a two whole screen triangles.
type GlWindow struct {
	width  int
	height int

	lastFrameTime time.Duration
	frameStart    time.Time
	frameTimeLock sync.RWMutex

	window *glfw.Window

	pixBufferLock sync.RWMutex
	pixBuffer     []float32
	pixSamples    []uint16

	glProgram uint32 // Holds the OpenGL program
	glVao     uint32 // Our only vertex array object
	glVbo     uint32 // The vertex buffer object which holds triangles and tex coords
	glEbo     uint32 // E(?) buffer object which enumerates the vbo vertices
	glTexture uint32 // Texture with the rendered scene

	glInited bool

	args GlWinArgs

	// Engine stuff

	sampler *sampler.SimpleSampler
	tracer  *engine.FPSEngine
	cam     camera.Camera
}

func NewGlWIndow(args GlWinArgs) *GlWindow {
	gwWin := &GlWindow{
		args: args,
	}
	return gwWin
}

func (g *GlWindow) Init(width int, height int) error {
	g.width = width
	g.height = height

	g.pixBuffer = make([]float32, g.width*g.height*3)
	g.pixSamples = make([]uint16, g.width*g.height)

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

	// Locking is slow, embrace the race!
	// g.pixBufferLock.RLock()
	// defer g.pixBufferLock.RUnlock()

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(g.width), int32(g.height),
		0, gl.RGB, gl.FLOAT, gl.Ptr(g.pixBuffer))
}

func (g *GlWindow) Wait() {
	gl.DeleteVertexArrays(1, &g.glVao)
	gl.DeleteBuffers(1, &g.glEbo)
	gl.DeleteBuffers(1, &g.glVbo)
}

func (g *GlWindow) StartFrame() {
	g.frameStart = time.Now()
	for i := 0; i < len(g.pixSamples); i++ {
		g.pixSamples[i] = 0
	}
}

func (g *GlWindow) DoneFrame() {
	g.frameTimeLock.Lock()
	defer g.frameTimeLock.Unlock()
	g.lastFrameTime = time.Since(g.frameStart)
}

func (g *GlWindow) LastFrameRederTime() time.Duration {
	g.frameTimeLock.RLock()
	defer g.frameTimeLock.RUnlock()
	return g.lastFrameTime
}

func (g *GlWindow) Set(x int, y int, clr color.Color) error {

	// Locking is slow, embrace the race!
	// g.pixBufferLock.Lock()
	// defer g.pixBufferLock.Unlock()

	sampleInd := g.width*y + x
	samples := g.pixSamples[sampleInd]

	oldWeight := float32(samples) / float32(samples+1)
	newWiehgt := 1 - oldWeight

	ri, gi, bi, _ := clr.RGBA()

	ind := g.width*y*3 + x*3
	g.pixBuffer[ind] = g.pixBuffer[ind]*oldWeight + newWiehgt*float32(ri)/65535.0
	g.pixBuffer[ind+1] = g.pixBuffer[ind+1]*oldWeight + newWiehgt*float32(gi)/65535.0
	g.pixBuffer[ind+2] = g.pixBuffer[ind+2]*oldWeight + newWiehgt*float32(bi)/65535.0

	g.pixSamples[sampleInd]++

	return nil
}

func (g *GlWindow) Width() int {
	return g.width
}

func (g *GlWindow) Height() int {
	return g.height
}

// Run initializes GLFW, creates window, initializes OpenGL and then starts the FPS
// engine tracing in this window.
func (g *GlWindow) Run() error {
	if err := g.initWindow(); err != nil {
		return fmt.Errorf("initGLFW: %w", err)
	}
	defer g.cleanWindow()

	winW, winH := g.window.GetFramebufferSize()
	if err := g.Init(winW, winH); err != nil {
		return fmt.Errorf("g.Init(): %w", err)
	}

	if err := g.initEngine(); err != nil {
		return fmt.Errorf("g.initEngine(): %w", err)
	}
	defer g.cleanEngine()

	g.mainLoop()

	return nil
}

func (g *GlWindow) initWindow() error {
	args := g.args

	if err := glfw.Init(); err != nil {
		return fmt.Errorf("initializing glfw failed. %w", err)
	}

	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	var err error
	var window *glfw.Window

	if args.Fullscreen {
		monitor := glfw.GetPrimaryMonitor()
		vm := monitor.GetVideoMode()
		monW, monH := vm.Width, vm.Height

		fmt.Printf("Running in fullscreen: %dx%d\n", monW, monH)

		window, err = glfw.CreateWindow(monW, monH, "Raytracer", monitor, nil)
	} else {
		window, err = glfw.CreateWindow(args.Width, args.Height, "Raytracer", nil, nil)
	}

	if err != nil {
		return fmt.Errorf("error creating window: %w", err)
	}

	window.MakeContextCurrent()
	g.window = window

	if args.VSync {
		window.MakeContextCurrent()
		glfw.SwapInterval(1)
	}

	window.SetCloseCallback(func(w *glfw.Window) {
		window.SetShouldClose(true)
	})

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int,
		action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape {
			window.SetShouldClose(true)
			return
		}
	})

	return nil
}

func (g *GlWindow) cleanWindow() {
	g.window.Destroy()
	glfw.Terminate()
}

func (g *GlWindow) initEngine() error {
	smpl := sampler.NewSimple(g.Width(), g.Height(), g)

	if g.args.Interactive {
		smpl.MakeContinuous()
	}

	cam := scene.GetCamera(float64(g.Width()), float64(g.Height()))

	tracer := engine.NewFPS(smpl)
	tracer.SetTarget(g, cam)
	tracer.ShowBBoxes = g.args.ShowBBoxes

	fmt.Printf("Loading scene...\n")
	loadingStart := time.Now()
	tracer.Scene.InitScene(g.args.SceneName)
	fmt.Printf("Loading scene took %s\n", time.Since(loadingStart))

	g.sampler = smpl
	g.tracer = tracer
	g.cam = cam

	return nil
}

func (g *GlWindow) cleanEngine() {
	g.sampler.Stop()
	g.tracer.StopRendering()
}

func (g *GlWindow) mainLoop() {
	g.tracer.Render()

	minFrameTime, _ := time.ParseDuration(
		fmt.Sprintf("%dms", int(1000.0/float32(g.args.FPSCap))),
	)

	g.window.MakeContextCurrent()

	var traceStarted bool
	var bPressed bool

	for !g.window.ShouldClose() {
		renderStart := time.Now()
		g.Render()
		renderTime := time.Since(renderStart)

		glfw.PollEvents()
		if g.args.Interactive {
			handleInteractionEvents(g.window, g.cam)

			if !bPressed && g.window.GetKey(glfw.KeyB) == glfw.Press {
				g.tracer.ShowBBoxes = !g.tracer.ShowBBoxes
				bPressed = true
			}

			if bPressed && g.window.GetKey(glfw.KeyB) == glfw.Release {
				bPressed = false
			}

			if !traceStarted && g.window.GetKey(glfw.KeyT) == glfw.Press {
				traceStarted = true
				go func() {
					collectTrace()
					traceStarted = false
				}()
			}
		}
		g.window.SwapBuffers()

		elapsed := time.Since(renderStart)
		if elapsed < minFrameTime {
			time.Sleep(minFrameTime - elapsed)
			elapsed = minFrameTime
		}

		if g.args.ShowFPS {
			fps := 1 / elapsed.Seconds()
			fmt.Printf("\r                                                               ")
			fmt.Printf("\rFPS: %5.3f Render time: %8s Last frame: %12s", fps, renderTime,
				g.LastFrameRederTime())
		}
	}

	fmt.Println("\nClosing window, rendering stopped.")
	g.Wait()
}

func handleInteractionEvents(window *glfw.Window, cam camera.Camera) {
	moveSpeed := 0.15
	rotateSpeed := 3.0
	if window.GetKey(glfw.KeyW) == glfw.Press {
		cam.Forward(moveSpeed)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		cam.Backward(moveSpeed)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		cam.Left(moveSpeed)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		cam.Right(moveSpeed)
	}
	if window.GetKey(glfw.KeyUp) == glfw.Press {
		cam.Pitch(rotateSpeed)
	}
	if window.GetKey(glfw.KeyDown) == glfw.Press {
		cam.Pitch(-rotateSpeed)
	}
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		cam.Yaw(-rotateSpeed)
	}
	if window.GetKey(glfw.KeyRight) == glfw.Press {
		cam.Yaw(rotateSpeed)
	}
}

func collectTrace() {
	traceFile := "trace.out"

	fh, err := os.Create(traceFile)

	if err != nil {
		fmt.Printf("Error creating trace file: %s\n", err)
		return
	}

	defer fh.Close()

	if err := trace.Start(fh); err != nil {
		fmt.Printf("Error staring trace: %s\n", err)
		return
	}

	defer trace.Stop()

	time.Sleep(2 * time.Second)
	fmt.Printf("Creating trace in %s\n", traceFile)
}
