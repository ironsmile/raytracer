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

const (
	WIDTH  = 1024
	HEIGHT = 768
)

var (
	cpuprofile  = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile  = flag.String("memprofile", "", "write memory profile to this file")
	filename    = flag.String("filename", "/tmp/rendered.png", "output file")
	interactive = flag.Bool("interactive", false, "starts the renderer in opengl win")
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
	err := output.Init(WIDTH, HEIGHT)
	if err != nil {
		log.Fatal("%s\n", err)
	}
	cam := MakePinholeCamera(output)
	tracer := engine.NewEngine()
	tracer.Scene.InitScene()
	tracer.SetTarget(output, cam)
	tracer.InitRender()
	renderTimer := time.Now()
	tracer.Render()
	fmt.Printf("Rendering finished: %s\n", time.Since(renderTimer))
	output.Wait()
}

func interactiveRenderer() {
	hasWindow := glfw3.Init()

	if !hasWindow {
		log.Fatal("Initializing glfw3 failed")
	}

	window, err := glfw3.CreateWindow(WIDTH, HEIGHT, "Raytracer", nil, nil)

	if err != nil {
		log.Fatal("%s\n", err.Error())
	}

	window.SetCloseCallback(func(w *glfw3.Window) {
		window.SetShouldClose(true)
	})

	output := film.NewGlWIndow(window)
	err = output.Init(WIDTH, HEIGHT)

	if err != nil {
		log.Fatal("%s\n", err.Error())
	}

	fmt.Println("Creating camera...")
	cam := MakePinholeCamera(output)

	window.SetKeyCallback(func(w *glfw3.Window, key glfw3.Key, scancode int,
		action glfw3.Action, mods glfw3.ModifierKey) {
		if key == glfw3.KeyEscape {
			window.SetShouldClose(true)
			return
		}
	})

	fmt.Println("Creating new engine...")
	tracer := engine.NewEngine()

	fmt.Println("Initializing scene...")
	tracer.Scene.InitScene()
	tracer.SetTarget(output, cam)

	fmt.Println("Initializing renderer...")
	tracer.InitRender()

	for !window.ShouldClose() {

		fmt.Println("Rendering...")
		renderTimer := time.Now()
		tracer.Render()
		fmt.Printf("Rendering finished: %s\n", time.Since(renderTimer))

		glfw3.WaitEvents()
		pollEvents(window, cam)
	}

	output.Wait()

	fmt.Println("Destroying window and terminating glfw3")
	window.Destroy()
	glfw3.Terminate()
}

func pollEvents(window *glfw3.Window, cam camera.Camera) {
	if window.GetKey(glfw3.KeyW) == glfw3.Press {
		cam.Forward(0.25)
	}
	if window.GetKey(glfw3.KeyS) == glfw3.Press {
		cam.Backward(0.25)
	}
	if window.GetKey(glfw3.KeyA) == glfw3.Press {
		cam.Left(0.25)
	}
	if window.GetKey(glfw3.KeyD) == glfw3.Press {
		cam.Right(0.25)
	}
}

func MakePinholeCamera(f film.Film) camera.Camera {
	pos := geometry.NewPoint(0, 0, -5)
	lookAtPoint := geometry.NewPoint(0, 0, 1)
	up := geometry.NewVector(0, 1, 0)

	return camera.NewPinholeCamera(pos, lookAtPoint, up, 1, f)
}
