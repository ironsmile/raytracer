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
	"slices"
	"strings"
	"time"

	"github.com/ironsmile/raytracer/engine"
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/sampler"
	"github.com/ironsmile/raytracer/scene"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// GLFW requires this.
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
	renderWidth = flag.Int("w", 1024,
		"image or window width in pixels")
	renderHeight = flag.Int("h", 768,
		"image or window height in pixels")
	fpsCap = flag.Uint("fps-cap", 30,
		"maximum number of frames per second. Zero means no FPS cap.")
	showBBoxes = flag.Bool("show-bboxes", false,
		"show bounding boxes around objects")
	withVulkan = flag.Bool("vulkan", false,
		"use Vulkan instead of OpenGL")
	sceneName = flag.String("scene", "teapot",
		"scene to render. Possible values: teapot, car")
	debugMode = flag.Bool("D", false,
		"debug mode, will print diagnostics information")
	debugRays = flag.String("debug-rays", "",
		"file nam which contains a list of rays which will be added to the scene\n"+
			"with deubgging purposes. For the format of the file see the code\n"+
			"comment on [scene.Scene.SetDebugRaysFile] function.")
)

func main() {

	flag.Parse()

	if !slices.Contains(scene.PossibleScenes, *sceneName) {
		log.Fatalf("scene must be one of: %s", strings.Join(scene.PossibleScenes, ", "))
	}

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

	if *debugRays != "" {
		scene.SetDebugRaysFile(*debugRays)
	}

	if *filename != "" {
		infileRenderer()
	} else if *withVulkan {
		vulkanWindowRenderer()
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
	if err := output.Init(*renderWidth, *renderHeight); err != nil {
		log.Fatalf("%s\n", err)
	}

	smpl := sampler.NewSimple(output.Width(), output.Height(), output)
	cam := scene.GetCamera(float64(output.Width()), float64(output.Height()))
	tracer := engine.New(smpl)
	tracer.SetTarget(output, cam)
	tracer.Scene.InitScene(*sceneName)
	tracer.ShowBBoxes = *showBBoxes

	renderTimer := time.Now()
	tracer.Render()
	fmt.Printf("Rendering finished: %s\n", time.Since(renderTimer))

	smpl.Stop()
	output.Wait()
}

func openglWindowRenderer() {
	args := film.GlWinArgs{
		Fullscreen:  *fullscreen,
		VSync:       *vsync,
		Width:       *renderWidth,
		Height:      *renderHeight,
		Interactive: *interactive,
		ShowBBoxes:  *showBBoxes,
		FPSCap:      *fpsCap,
		ShowFPS:     *showFPS,
		SceneName:   *sceneName,
	}
	glWin := film.NewGlWIndow(args)
	if err := glWin.Run(); err != nil {
		log.Fatalf("Error running GL window: %s", err)
	}
}

func vulkanWindowRenderer() {
	args := film.VulkanAppArgs{
		Debug:       *debugMode,
		Fullscreen:  *fullscreen,
		VSync:       *vsync,
		Width:       *renderWidth,
		Height:      *renderHeight,
		Interactive: *interactive,
		ShowBBoxes:  *showBBoxes,
		FPSCap:      *fpsCap,
		ShowFPS:     *showFPS,
		SceneName:   *sceneName,
	}

	app := film.NewVulkanWindow(args)
	if err := app.Run(); err != nil {
		log.Fatalf("error running: %s", err)
	}
}
