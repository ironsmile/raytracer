package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ironsmile/raytracer/camera"
	"github.com/ironsmile/raytracer/engine"
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/geometry"

	"github.com/go-gl/glfw3"
)

const (
	WIDTH  = 1024
	HEIGHT = 768
)

func main() {

	hasWindow := glfw3.Init()

	if !hasWindow {
		fmt.Errorf("Initializing glfw3 failed")
		os.Exit(1)
	}

	window, err := glfw3.CreateWindow(WIDTH, HEIGHT, "Raytracer", nil, nil)

	if err != nil {
		fmt.Errorf("%s\n", err.Error())
		os.Exit(1)
	}

	window.SetCloseCallback(func(w *glfw3.Window) {
		window.SetShouldClose(true)
	})

	output := film.NewGlWIndow(window)
	// output := film.NewNullFilm()
	// output := film.NewImage("/tmp/rendered.png")
	err = output.Init(WIDTH, HEIGHT)

	if err != nil {
		fmt.Errorf("%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Creating camera...")
	// cam := MakeStackOverflowCamera(output)
	cam := MakePinholeCamera(output)
	// cam := MakePerspectiveCamera(output)

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
		_ = tracer.Render()
		fmt.Printf("Rendering finished - %s\n", time.Since(renderTimer))

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

// func MakeStackOverflowCamera(f film.Film) camera.Camera {
// 	pos := geometry.NewPoint(0, 0, -5)
// 	lookAtPoint := geometry.NewPoint(0, 0, 1)
// 	up := geometry.NewVector(0, 1, 0)

// 	return camera.NewStackOverflowCamera(pos, lookAtPoint, up, f)
// }
