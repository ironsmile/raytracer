package unsafer

import (
	"reflect"
	"unsafe"
)

// StructToBytes returns a byte slice which points to the memory of input.
//
// Note that this function does not create a copy. It points to the exact
// same memory as the input itself.
func StructToBytes[T any](input *T) []byte {
	inputSize := int(unsafe.Sizeof(*input))

	header := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(input)),
		Len:  inputSize,
		Cap:  inputSize,
	}

	return *(*[]byte)(unsafe.Pointer(&header))
}
