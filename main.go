package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ironsmile/raytracer/engine"
	"github.com/ironsmile/raytracer/film"
)

const (
	WIDTH  = 1024
	HEIGHT = 768
)

func main() {
	win := film.NewGlWIndow()
	// win := film.NewNullFilm()
	// win := film.NewImage("/tmp/rendered.png")
	err := win.Init(WIDTH, HEIGHT)

	if err != nil {
		fmt.Errorf("%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Creating new engine...")
	tracer := engine.NewEngine()

	fmt.Println("Initializing scene...")
	tracer.Scene.InitScene()
	tracer.SetTarget(win)

	fmt.Println("Initializing renderer...")
	tracer.InitRender()

	renderTimer := time.Now()
	fmt.Println("Rendering...")
	_ = tracer.Render()
	fmt.Printf("Rendering finished - %s\n", time.Since(renderTimer))

	win.Wait()
}
