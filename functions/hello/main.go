// Package main : this is a wasm module exposing a hello function
package main

import (
	"unsafe"
)


// Read memory
func readMemory(bufferPosition *uint32, length uint32) []byte {
	subjectBuffer := make([]byte, length)
	pointer := uintptr(unsafe.Pointer(bufferPosition))
	for i := 0; i < int(length); i++ {
		s := *(*int32)(unsafe.Pointer(pointer + uintptr(i)))
		subjectBuffer[i] = byte(s)
	}
	return subjectBuffer
}


func copyToMemory(buffer []byte) (uint32, uint32) {
	bufferPtr := &buffer[0]
	unsafePtr := uintptr(unsafe.Pointer(bufferPtr))

	pos := uint32(unsafePtr)
	size := uint32(len(buffer))

	return pos, size
}

func pack(pos uint32, size uint32) uint64 {
	return (uint64(pos) << uint64(32)) | uint64(size)
}


// hello function
//
//export hello
func hello(valuePosition *uint32, length uint32) uint64 {

	// read the memory to get the argument(s)
	valueBytes := readMemory(valuePosition, length)

	message := "👋 Hello " + string(valueBytes) + " 😃"

	// copy the value to memory
	// get the position and the size of the buffer (in memory)
	pos, size := copyToMemory([]byte(message))

	// return the position and size
	return pack(pos, size)

}

func main() {}
