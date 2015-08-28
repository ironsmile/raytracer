package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/ironsmile/raytracer/camera"
	"github.com/ironsmile/raytracer/engine"
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/sampler"
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
	if err := output.Init(*WIDTH, *HEIGHT); err != nil {
		log.Fatal("%s\n", err)
	}

	smpl := &sampler.SimpleSampler{}

	if err := smpl.Init(output); err != nil {
		log.Fatal("%s\n", err)
	}

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

func interactiveRenderer() {

	if err := glfw.Init(); err != nil {
		log.Fatal("Initializing glfw failed. %s", err)
	}
	defer glfw.Terminate()

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

	// fmt.Printf("swap interval: %t\n", glfw.ExtensionSupported("SwapInterval"))

	if err != nil {
		log.Fatal("%s\n", err.Error())
	}

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

	smpl := &sampler.SimpleSampler{}
	smpl.MakeContinuous()

	if err = smpl.Init(output); err != nil {
		log.Fatal("%s\n", err)
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
	tracer.Scene.InitScene()

	// window.MakeContextCurrent()
	// glfw.SwapInterval(1)

	tracer.Render()

	for !window.ShouldClose() {
		// glfw.WaitEvents()
		time.Sleep(25 * time.Millisecond)
		glfw.PollEvents()
		pollEvents(window, cam)
	}

	smpl.Stop()
	tracer.StopRendering()
}

func pollEvents(window *glfw.Window, cam camera.Camera) {
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

	return camera.NewPinholeCamera(pos, lookAtPoint, up, 1, f)
}
