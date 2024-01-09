package main

import (
	"testing"

	"github.com/ironsmile/raytracer/engine"
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/sampler"
	"github.com/ironsmile/raytracer/scene"
)

// The benchmark result without the sampler is 806732319 ns/op
// Latest benchmark for the sampler            1130174954 ns/op
func BenchmarkImageCreation(t *testing.B) {
	output := film.NewImage("/dev/null")
	if err := output.Init(1024, 768); err != nil {
		t.Fatalf("Initializing nil output failed. %s", err)
	}
	smpl := sampler.NewSimple(output.Width(), output.Height(), output)
	cam := scene.GetCamera(float64(output.Width()), float64(output.Height()))
	tracer := engine.New(smpl)
	tracer.SetTarget(output, cam)
	tracer.Scene.InitScene("teapot")

	for i := 0; i < t.N; i++ {
		output.StartFrame()
		tracer.Render()
		output.DoneFrame()
	}

	smpl.Stop()
	output.Wait()
}
