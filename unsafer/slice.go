package unsafer

import (
	"reflect"
	"unsafe"

	vk "github.com/vulkan-go/vulkan"
)

// SliceToBytes interprets an arbitrary input slice as a byte slice.
//
// Note that the returned slice points to the same underlying data in memory. It
// does not make a copy.
func SliceToBytes[T any](input []T) []byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&input))
	header.Len = int(unsafe.Sizeof(input[0])) * len(input)
	header.Cap = header.Len
	bytesSlice := *(*[]byte)(unsafe.Pointer(&header))
	return bytesSlice
}

// SliceBytesToUint32 copies the slice of bytes into a slice of uint32. Note that
// data's length must be divisible by four since uint32 is four bytes. It does that
// with memcopy.
func SliceBytesToUint32(data []byte) []uint32 {
	buf := make([]uint32, len(data)/4)
	vk.Memcopy(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&buf)).Data), data)
	return buf
}
