package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/go-gl/glfw3"

	"github.com/ironsmile/raytracer/camera"
	"github.com/ironsmile/raytracer/engine"
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/geometry"
)

var (
	cpuprofile = flag.String("cpuprofile", "",
		"write cpu profile to file")
	memprofile = flag.String("memprofile", "",
		"write memory profile to this file")
	filename = flag.String("filename", "/tmp/rendered.png",
		"output file")
	interactive = flag.Bool("interactive", false,
		"starts the renderer in opengl win")
	vsync = flag.Bool("vsync", true,
		"control vsync for interactive renderer")
	fullscreen = flag.Bool("fullscreen", false,
		"run fullscreen in native resolution")
	WIDTH = flag.Int("w", 1024,
		"image or window width in pixels")
	HEIGHT = flag.Int("h", 768,
		"image or window height in pixels")
)

func main() {

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *interactive {
		interactiveRenderer()
	} else {
		infileRenderer()
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
	err := output.Init(*WIDTH, *HEIGHT)
	if err != nil {
		log.Fatal("%s\n", err)
	}
	cam := MakePinholeCamera(output)
	tracer := engine.NewEngine()
	tracer.SetTarget(output, cam)
	tracer.InitRender()
	tracer.Scene.InitScene()
	renderTimer := time.Now()
	tracer.Render()
	fmt.Printf("Rendering finished: %s\n", time.Since(renderTimer))
	output.Wait()
}

func interactiveRenderer() {

	if !glfw3.Init() {
		log.Fatal("Initializing glfw3 failed")
	}
	defer glfw3.Terminate()

	var err error
	var window *glfw3.Window

	if *fullscreen {
		monitor, err := glfw3.GetPrimaryMonitor()

		if err != nil {
			log.Fatal("Was not able to find primary monitor")
		}

		vm, err := monitor.GetVideoMode()

		if err != nil {
			log.Fatal("Was not able to get monitor video mode")
		}
		monW, monH := vm.Width, vm.Height

		fmt.Printf("Running in fullscreen: %dx%d\n", monW, monH)

		window, err = glfw3.CreateWindow(monW, monH, "Raytracer", monitor, nil)
	} else {
		window, err = glfw3.CreateWindow(*WIDTH, *HEIGHT, "Raytracer", nil, nil)
	}

	// fmt.Printf("swap interval: %t\n", glfw3.ExtensionSupported("SwapInterval"))

	if err != nil {
		log.Fatal("%s\n", err.Error())
	}

	defer func() {
		window.MakeContextCurrent()
		window.Destroy()
	}()

	if *vsync {
		window.MakeContextCurrent()
		glfw3.SwapInterval(1)
	}

	window.SetCloseCallback(func(w *glfw3.Window) {
		window.SetShouldClose(true)
	})

	output := film.NewGlWIndow(window)
	winW, winH := window.GetFramebufferSize()
	err = output.Init(winW, winH)

	if err != nil {
		log.Fatal("%s\n", err.Error())
	}

	cam := MakePinholeCamera(output)

	window.SetKeyCallback(func(w *glfw3.Window, key glfw3.Key, scancode int,
		action glfw3.Action, mods glfw3.ModifierKey) {
		if key == glfw3.KeyEscape {
			window.SetShouldClose(true)
			return
		}
	})

	tracer := engine.NewFPSEngine()
	tracer.SetTarget(output, cam)
	tracer.InitRender()
	tracer.Scene.InitScene()

	// window.MakeContextCurrent()
	// glfw3.SwapInterval(1)

	tracer.Render()

	for !window.ShouldClose() {
		// glfw3.WaitEvents()
		time.Sleep(25 * time.Millisecond)
		glfw3.PollEvents()
		pollEvents(window, cam)
	}

	tracer.StopRendering()
}

func pollEvents(window *glfw3.Window, cam camera.Camera) {
	moveSpeed := 0.15
	rotateSpeed := 3.0
	if window.GetKey(glfw3.KeyW) == glfw3.Press {
		cam.Forward(moveSpeed)
	}
	if window.GetKey(glfw3.KeyS) == glfw3.Press {
		cam.Backward(moveSpeed)
	}
	if window.GetKey(glfw3.KeyA) == glfw3.Press {
		cam.Left(moveSpeed)
	}
	if window.GetKey(glfw3.KeyD) == glfw3.Press {
		cam.Right(moveSpeed)
	}
	if window.GetKey(glfw3.KeyUp) == glfw3.Press {
		cam.Pitch(rotateSpeed)
	}
	if window.GetKey(glfw3.KeyDown) == glfw3.Press {
		cam.Pitch(-rotateSpeed)
	}
	if window.GetKey(glfw3.KeyLeft) == glfw3.Press {
		cam.Yaw(-rotateSpeed)
	}
	if window.GetKey(glfw3.KeyRight) == glfw3.Press {
		cam.Yaw(rotateSpeed)
	}
}

func MakePinholeCamera(f film.Film) camera.Camera {
	pos := geometry.NewPoint(0, 0, -5)
	lookAtPoint := geometry.NewPoint(0, 0, 1)
	up := geometry.NewVector(0, 1, 0)

	return camera.NewPinholeCamera(pos, lookAtPoint, up, 1, f)
}
