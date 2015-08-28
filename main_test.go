package main

import (
	"testing"

	"github.com/ironsmile/raytracer/engine"
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/sampler"
)

// The benchmark result without the sampler is 806732319 ns/op
// Latest benchmark for the sampler            1130174954 ns/op
func BenchmarkImageCreation(t *testing.B) {
	output := film.NewImage("/dev/null")
	if err := output.Init(1024, 768); err != nil {
		t.Fatalf("Initializing nil output failed. %s", err)
	}
	smpl := &sampler.SimpleSampler{}

	if err := smpl.Init(output); err != nil {
		t.Fatalf("Initializing sampler failed. %s", err)
	}

	cam := MakePinholeCamera(output)
	tracer := engine.New(smpl)
	tracer.SetTarget(output, cam)
	tracer.Scene.InitScene()

	for i := 0; i < t.N; i++ {
		output.StartFrame()
		tracer.Render()
		output.DoneFrame()
	}

	smpl.Stop()
	output.Wait()
}
