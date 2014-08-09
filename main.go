package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ironsmile/raytracer/camera"
	"github.com/ironsmile/raytracer/engine"
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/geometry"
)

const (
	WIDTH  = 1024
	HEIGHT = 768
)

func main() {
	output := film.NewGlWIndow()
	// output := film.NewNullFilm()
	// output := film.NewImage("/tmp/rendered.png")
	err := output.Init(WIDTH, HEIGHT)

	if err != nil {
		fmt.Errorf("%s\n", err.Error())
		os.Exit(1)
	}

	cam := camera.NewBasicCamera()
	cam.SetOrigin(geometry.NewPoint(0, 0, -5))

	fmt.Println("Creating new engine...")
	tracer := engine.NewEngine()

	fmt.Println("Initializing scene...")
	tracer.Scene.InitScene()
	tracer.SetTarget(output, cam)

	fmt.Println("Initializing renderer...")
	tracer.InitRender()

	fmt.Println("Rendering...")
	renderTimer := time.Now()
	_ = tracer.Render()
	fmt.Printf("Rendering finished - %s\n", time.Since(renderTimer))

	output.Wait()
}
