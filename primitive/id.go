package primitive

import (
	"sync/atomic"
)

var nextID uint64

// GetNewID returns a new unique ID
func GetNewID() uint64 {
	return atomic.AddUint64(&nextID, 1) - 1
}
