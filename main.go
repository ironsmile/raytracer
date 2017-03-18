package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/ironsmile/raytracer/camera"
	"github.com/ironsmile/raytracer/engine"
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/sampler"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// GLFW3.1 requires this.
	runtime.LockOSThread()
}

var (
	cpuprofile = flag.String("cpuprofile", "",
		"write cpu profile to file")
	memprofile = flag.String("memprofile", "",
		"write memory profile to this file")
	filename = flag.String("filename", "",
		"output the image to a PNG file instead of showing it to the screen")
	interactive = flag.Bool("interactive", false,
		"starts the renderer in interactive mode")
	vsync = flag.Bool("vsync", true,
		"control vsync for interactive renderer")
	showFPS = flag.Bool("show-fps", true,
		"continuously print the OpenGL FPS stats in the console")
	fullscreen = flag.Bool("fullscreen", false,
		"run fullscreen in native resolution")
	WIDTH = flag.Int("w", 1024,
		"image or window width in pixels")
	HEIGHT = flag.Int("h", 768,
		"image or window height in pixels")
	fpsCap = flag.Uint("fps-cap", 30,
		"maximum number of frames per second")
)

func main() {

	flag.Parse()

	go func() {
		log.Println(http.ListenAndServe("localhost:6464", nil))
	}()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *filename != "" {
		infileRenderer()
	} else {
		openglWindowRenderer()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
}

func infileRenderer() {
	output := film.NewImage(*filename)
	if err := output.Init(*WIDTH, *HEIGHT); err != nil {
		log.Fatalf("%s\n", err)
	}

	smpl := sampler.NewSimple(output)
	cam := MakePinholeCamera(output)
	tracer := engine.New(smpl)
	tracer.SetTarget(output, cam)
	tracer.Scene.InitScene()

	renderTimer := time.Now()
	tracer.Render()
	fmt.Printf("Rendering finished: %s\n", time.Since(renderTimer))

	smpl.Stop()
	output.Wait()
}

func openglWindowRenderer() {

	if err := glfw.Init(); err != nil {
		log.Fatal("Initializing glfw failed. %s", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	var err error
	var window *glfw.Window

	if *fullscreen {
		monitor := glfw.GetPrimaryMonitor()
		vm := monitor.GetVideoMode()
		monW, monH := vm.Width, vm.Height

		fmt.Printf("Running in fullscreen: %dx%d\n", monW, monH)

		window, err = glfw.CreateWindow(monW, monH, "Raytracer", monitor, nil)
	} else {
		window, err = glfw.CreateWindow(*WIDTH, *HEIGHT, "Raytracer", nil, nil)
	}

	if err != nil {
		log.Fatal("%s\n", err.Error())
	}

	window.MakeContextCurrent()

	defer func() {
		window.MakeContextCurrent()
		window.Destroy()
	}()

	if *vsync {
		window.MakeContextCurrent()
		glfw.SwapInterval(1)
	}

	window.SetCloseCallback(func(w *glfw.Window) {
		window.SetShouldClose(true)
	})

	output := film.NewGlWIndow(window)
	winW, winH := window.GetFramebufferSize()
	if err := output.Init(winW, winH); err != nil {
		log.Fatal("%s\n", err.Error())
	}

	smpl := sampler.NewSimple(output)

	if *interactive {
		smpl.MakeContinuous()
	}

	cam := MakePinholeCamera(output)

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int,
		action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape {
			window.SetShouldClose(true)
			return
		}
	})

	tracer := engine.NewFPS(smpl)
	tracer.SetTarget(output, cam)

	fmt.Printf("Loading scene...\n")
	loadingStart := time.Now()
	tracer.Scene.InitScene()
	fmt.Printf("Loading scene took %s\n", time.Since(loadingStart))

	tracer.Render()

	minFrameTime, _ := time.ParseDuration(fmt.Sprintf("%dms", int(1000.0/float32(*fpsCap))))

	window.MakeContextCurrent()
	for !window.ShouldClose() {
		renderStart := time.Now()
		output.Render()
		renderTime := time.Since(renderStart)

		glfw.PollEvents()
		if *interactive {
			handleInteractionEvents(window, cam)
		}
		window.SwapBuffers()

		elapsed := time.Since(renderStart)
		if elapsed < minFrameTime {
			time.Sleep(minFrameTime - elapsed)
			elapsed = minFrameTime
		}

		if *showFPS {
			fps := 1 / elapsed.Seconds()
			fmt.Printf("\r                                                               ")
			fmt.Printf("\rFPS: %5.3f Render time: %8s Last frame: %12s", fps, renderTime,
				output.LastFrameRederTime())
		}
	}

	fmt.Println("\nClosing window, rendering stopped.")
	output.Wait()

	smpl.Stop()
	tracer.StopRendering()
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

func MakePinholeCamera(f film.Film) camera.Camera {
	pos := geometry.NewPoint(0, 0, -5)
	lookAtPoint := geometry.NewPoint(0, 0, 1)
	up := geometry.NewVector(0, 1, 0)

	return camera.NewPinhole(pos, lookAtPoint, up, 1, f)
}
