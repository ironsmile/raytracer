package film

import (
	"image/color"
	"sync"
	"time"

	vk "github.com/vulkan-go/vulkan"
)

type vulkanFilm struct {
	pixBuffer  []uint8
	pixSamples []uint16

	pixBufferFormat vk.Format

	width  uint32
	height uint32

	lastFrameTime time.Duration
	frameStart    time.Time
	frameTimeLock *sync.RWMutex
}

func newVulkanFilm(width, height uint32) *vulkanFilm {
	return &vulkanFilm{
		width:           width,
		height:          height,
		pixBuffer:       make([]uint8, width*height*4),
		pixSamples:      make([]uint16, width*height),
		pixBufferFormat: vk.FormatR8g8b8a8Srgb,

		frameTimeLock: &sync.RWMutex{},
	}
}

func (f *vulkanFilm) Set(x int, y int, clr color.Color) error {
	sampleInd := f.width*uint32(y) + uint32(x)
	samples := f.pixSamples[sampleInd]

	oldWeight := float32(samples) / float32(samples+1)
	newWeight := 1 - oldWeight

	ri, gi, bi, _ := clr.RGBA()

	nc := (newWeight * 255) / 0xffff

	ind := f.width*uint32(y)*4 + uint32(x)*4
	f.pixBuffer[ind] = uint8(float32(f.pixBuffer[ind])*oldWeight + float32(ri)*nc)
	f.pixBuffer[ind+1] = uint8(float32(f.pixBuffer[ind+1])*oldWeight + float32(gi)*nc)
	f.pixBuffer[ind+2] = uint8(float32(f.pixBuffer[ind+2])*oldWeight + float32(bi)*nc)

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

// asVkBuffer interprets returns the film pixel buffer as a byte slice suitable
// for copying in a buffer for a Vulkan Image.
//
// This function uses the pixel buffer data "in-place" where possible in order
// to avoid copying.
func (f *vulkanFilm) asVkBuffer() []byte {
	return f.pixBuffer
}

func (f *vulkanFilm) getFormat() vk.Format {
	return f.pixBufferFormat
}

// getBufferSize returns the pixel buffer size in bytes.
func (f *vulkanFilm) getBufferSize() uint64 {
	return uint64(f.width) * uint64(f.height) * 4
}
