package main

import (
	"testing"

	"github.com/ironsmile/raytracer/engine"
	"github.com/ironsmile/raytracer/film"
	"github.com/ironsmile/raytracer/sampler"
)

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
	tracer := engine.NewEngine()
	tracer.SetTarget(output, cam)
	tracer.InitRender(smpl)
	tracer.Scene.InitScene()

	for i := 0; i < t.N; i++ {
		tracer.Render()
	}

	smpl.Stop()
	output.Wait()
}
