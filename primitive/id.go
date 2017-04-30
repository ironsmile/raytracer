package primitive

import (
	"sync/atomic"
)

const unnamedPrimitive = "Unnamed Primitive"

var primNames map[uint64]string
var nextID uint64

// GetNewID returns a new unique ID
func GetNewID() uint64 {
	return atomic.AddUint64(&nextID, 1) - 1
}

// SetName sets a global name for a primitive
func SetName(id uint64, name string) {
	if primNames == nil {
		primNames = make(map[uint64]string)
	}
	primNames[id] = name
}

// GetName returns the global name of a primitive by its ID
func GetName(id uint64) string {
	if primNames == nil {
		return unnamedPrimitive
	}

	if name, ok := primNames[id]; ok {
		return name
	}

	return unnamedPrimitive
}
