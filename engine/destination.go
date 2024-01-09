package engine

type Destination interface {
	StartFrame()
	DoneFrame()
	Wait()
	Width() int
	Height() int
}
