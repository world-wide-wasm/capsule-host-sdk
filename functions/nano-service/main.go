// Package main : this is a wasm module exposing a hello function
package main

import (
	"encoding/json"
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

//export hostPrintStr
func hostPrintStr(messagePosition, messageLength uint32) uint32

// printStr calls a host function: hostPrintStr
// and prints a string
func printStr(message string) {
	messagePosition, messageSize := copyToMemory([]byte(message))
	hostPrintStr(messagePosition, messageSize)
}

// functionHandler function
//
//export functionHandler
func functionHandler(valuePosition *uint32, length uint32) uint64 {

	// read the memory to get the argument(s)
	valueBytes := readMemory(valuePosition, length)

	printStr("🤖 parameter: " + string(valueBytes))

	//jsonStr := string(valueBytes)

	var data map[string]interface{}

	//err := json.Unmarshal([]byte(jsonStr), &data)

	err := json.Unmarshal(valueBytes, &data)
	if err != nil {
		printStr("😡 Error: " + err.Error())
	}

	printStr("📝: " + data["type"].(string))
	printStr("👤: " + data["body"].(string))

	if data["type"].(string) == "json" {
		var bodyMap map[string]interface{}
		//err := json.Unmarshal([]byte(jsonStr), &data)
		err := json.Unmarshal([]byte(data["body"].(string)), &bodyMap)

		if err != nil {
			printStr("😡 Error: " + err.Error())
		}

		for key, element := range bodyMap {
			printStr("Key: " + key + " => " + element.(string))
		}

	}

	responseData := map[string]interface{}{
		"body": `{"message": "👋 Hello Peeps🤗"}`,
		"header": map[string][]string{
			"Content-Type": {"application/json; charset=utf-8"},
		},
	}
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		pos, size := copyToMemory([]byte(`{"error":"`+err.Error()+`"}`))
		return pack(pos, size)
	}

	// copy the value to memory
	// get the position and the size of the buffer (in memory)
	pos, size := copyToMemory(jsonData)

	// return the position and size
	return pack(pos, size)

}

func main() {}
