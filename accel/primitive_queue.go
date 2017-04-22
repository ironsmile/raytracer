package accel

import (
	"container/heap"

	"github.com/ironsmile/raytracer/primitive"
)

// Item represents an item in the PriorityQueue
type Item struct {
	value primitive.Primitive
	index int
}

// PrimitiveQueue implements a queue of primitives
type PrimitiveQueue []*Item

// Len is needed to implement the heap interface
func (pq PrimitiveQueue) Len() int {
	return len(pq)
}

// Less is needed to implement the heap interface
func (pq PrimitiveQueue) Less(i, j int) bool {
	return pq[i].value.GetID() < pq[j].value.GetID()
}

func (pq PrimitiveQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push is needed to implement the heap interface
func (pq *PrimitiveQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop is needed to implement the heap interface
func (pq *PrimitiveQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// PushPrimitive is a convinience wrapper around heap.Push
func (pq *PrimitiveQueue) PushPrimitive(x primitive.Primitive) {
	heap.Push(pq, x)
}

// PopPrimitive is a convinience function for Popping from the queue and unwrapping
// the popped value from interface{} to the concrete type - Primitive
func (pq *PrimitiveQueue) PopPrimitive() primitive.Primitive {
	item := heap.Pop(pq).(*Item)
	return item.value
}

// Init is a convinience wrapper around heap.Init
func (pq *PrimitiveQueue) Init(x primitive.Primitive) {
	heap.Init(pq)
}

// NewPriorityQueue returns a new and initialized queue of primitives
func NewPriorityQueue(p []primitive.Primitive) *PrimitiveQueue {
	pq := make(PrimitiveQueue, len(p))

	for i, prim := range p {
		pq[i] = &Item{
			value: prim,
			index: i,
		}
	}

	heap.Init(&pq)

	return &pq
}
