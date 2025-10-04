package film

import (
	"image/color"
	"sync"
	"time"
)

type vulkanFilm struct {
	pixBuffer  []float32
	pixSamples []uint16

	width  uint32
	height uint32

	lastFrameTime time.Duration
	frameStart    time.Time
	frameTimeLock *sync.RWMutex
}

func newVulkanFilm(width, height uint32) *vulkanFilm {
	return &vulkanFilm{
		width:      width,
		height:     height,
		pixBuffer:  make([]float32, width*height*3),
		pixSamples: make([]uint16, width*height),

		frameTimeLock: &sync.RWMutex{},
	}
}

func (f *vulkanFilm) Set(x int, y int, clr color.Color) error {
	sampleInd := f.width*uint32(y) + uint32(x)
	samples := f.pixSamples[sampleInd]

	oldWeight := float32(samples) / float32(samples+1)
	newWeight := 1 - oldWeight

	ri, gi, bi, _ := clr.RGBA()

	ind := f.width*uint32(y)*3 + uint32(x)*3
	f.pixBuffer[ind] = f.pixBuffer[ind]*oldWeight + (float32(ri)/0xffff)*newWeight
	f.pixBuffer[ind+1] = f.pixBuffer[ind+1]*oldWeight + (float32(gi)/0xffff)*newWeight
	f.pixBuffer[ind+2] = f.pixBuffer[ind+2]*oldWeight + (float32(bi)/0xffff)*newWeight

	f.pixSamples[sampleInd]++

	return nil
}

func (f *vulkanFilm) DoneFrame() {
	f.frameTimeLock.Lock()
	defer f.frameTimeLock.Unlock()
	f.lastFrameTime = time.Since(f.frameStart)
}

func (f *vulkanFilm) StartFrame() {
	f.frameStart = time.Now()
	for i := 0; i < len(f.pixSamples); i++ {
		f.pixSamples[i] = 0
	}
}

func (f *vulkanFilm) FrameTime() time.Duration {
	f.frameTimeLock.Lock()
	defer f.frameTimeLock.Unlock()
	return f.lastFrameTime
}

func (f *vulkanFilm) Width() int {
	return int(f.width)
}

func (f *vulkanFilm) Height() int {
	return int(f.height)
}

func (f *vulkanFilm) Wait() {

}
