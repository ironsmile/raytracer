package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ironsmile/raytracer/camera"
	"github.com/ironsmile/raytracer/engine"
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/transform"
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

	fmt.Println("Creating camera...")
	cam := MakePinholeCamera(output)
	// cam := MakePerspectiveCamera(output)
	// cam := camera.NewDoychoCamera(geometry.NewPoint(0, 0, -5))

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

func MakePerspectiveCamera(f film.Film) camera.Camera {
	camOrigin := geometry.NewPoint(0, 0, -5)

	camToWorld := transform.LookAt(camOrigin,
		geometry.NewPoint(0, 0, -1),
		geometry.NewVector(0, 1, 0))

	oStart := camToWorld.Point(geometry.NewPoint(0, 0, 0))
	fmt.Printf("Cam Origin: %s\n", oStart)

	sOpen := 0.0
	sClose := 1.0
	lenRad := 1.5
	focalDist := 1e30
	frame := float64(WIDTH) / float64(HEIGHT)
	fov := 90.0

	screen := [4]float64{}

	if frame > 1.0 {
		screen[0] = -frame
		screen[1] = frame
		screen[2] = -1.0
		screen[3] = 1.0
	} else {
		screen[0] = -1.0
		screen[1] = 1.0
		screen[2] = -1.0 / frame
		screen[3] = 1.0 / frame
	}

	fmt.Printf("cameraToWorld transformation:\n%s\n", camToWorld)

	return camera.NewPerspectiveCamera(camToWorld, screen, sOpen, sClose, lenRad,
		focalDist, fov, f)
}

func MakePinholeCamera(f film.Film) camera.Camera {
	camToWorld := transform.LookAt(
		geometry.NewPoint(0, 0, -5),
		geometry.NewPoint(0, 0, -1),
		geometry.NewVector(0, 1, 0))

	return camera.NewPinholeCamera(camToWorld, 1e-3, f)
}
