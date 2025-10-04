package film

import (
	"fmt"
	"os"
	"runtime/trace"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ironsmile/raytracer/camera"
)

var (
	moveSpeed   = 3.0
	rotateSpeed = 25.0
)

func handleInteractionEvents(window *glfw.Window, cam camera.Camera, dur time.Duration) {
	// camera transform controls.
	if window.GetKey(glfw.KeyW) == glfw.Press {
		cam.Forward(moveSpeed * dur.Seconds())
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		cam.Backward(moveSpeed * dur.Seconds())
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		cam.Left(moveSpeed * dur.Seconds())
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		cam.Right(moveSpeed * dur.Seconds())
	}
	if window.GetKey(glfw.KeyE) == glfw.Press {
		cam.Up(moveSpeed * dur.Seconds())
	}
	if window.GetKey(glfw.KeyQ) == glfw.Press {
		cam.Down(moveSpeed * dur.Seconds())
	}
	if window.GetKey(glfw.KeyUp) == glfw.Press {
		cam.Pitch(rotateSpeed * dur.Seconds())
	}
	if window.GetKey(glfw.KeyDown) == glfw.Press {
		cam.Pitch(-rotateSpeed * dur.Seconds())
	}
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		cam.Yaw(-rotateSpeed * dur.Seconds())
	}
	if window.GetKey(glfw.KeyRight) == glfw.Press {
		cam.Yaw(rotateSpeed * dur.Seconds())
	}

	// movement speed controls.
	if window.GetKey(glfw.Key1) == glfw.Press {
		moveSpeed = 0.5
		rotateSpeed = 12
	}
	if window.GetKey(glfw.Key2) == glfw.Press {
		moveSpeed = 1
		rotateSpeed = 18
	}
	if window.GetKey(glfw.Key3) == glfw.Press {
		moveSpeed = 3
		rotateSpeed = 25
	}
	if window.GetKey(glfw.Key4) == glfw.Press {
		moveSpeed = 8
		rotateSpeed = 35
	}
	if window.GetKey(glfw.Key5) == glfw.Press {
		moveSpeed = 16
		rotateSpeed = 50
	}
}

func collectTrace() {
	traceFile := "trace.out"

	fh, err := os.Create(traceFile)

	if err != nil {
		fmt.Printf("Error creating trace file: %s\n", err)
		return
	}

	defer fh.Close()

	if err := trace.Start(fh); err != nil {
		fmt.Printf("Error staring trace: %s\n", err)
		return
	}

	defer trace.Stop()

	time.Sleep(2 * time.Second)
	fmt.Printf("Creating trace in %s\n", traceFile)
}
