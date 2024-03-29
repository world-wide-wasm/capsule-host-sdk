// Package wasm :
package wasm

import (
	"context"
	"errors"

	"github.com/tetratelabs/wazero/api"
)

/*
CallModuleFunction executes the wasm function of a module
*/
func CallModuleFunction(module api.Module, wasmFunctionName string, param []byte, ctx context.Context) ([]byte, error) {
	// These function are exported by TinyGo
	malloc := module.ExportedFunction("malloc")
	free := module.ExportedFunction("free")

	// 4- Get the reference to the Wasm function:
	wasmFunction := module.ExportedFunction(wasmFunctionName)

	paramSize := uint64(len(param))
	// Allocate Memory for param
	results, err := malloc.Call(ctx, paramSize)
	if err != nil {
		return nil, err
	}
	paramPosition := results[0]

	// Free the pointer when finished
	defer free.Call(ctx, paramPosition)

	/*
	allocate := module.ExportedFunction("allocate")
	deallocate := module.ExportedFunction("deallocate")
	*/

	// Copy param value to memory
	success := module.Memory().Write(uint32(paramPosition), param)
	if !success {
		return nil, errors.New("out of range of memory size")
	}

	// 6- Call function(pos, size)
	// Call the function with the position and the size of the value of param
	// The result type is []uint64
	result, err := wasmFunction.Call(ctx, paramPosition, paramSize)
	if err != nil {
		return nil, err
	}
	// Extract the position and size of from result
	valuePosition := uint32(result[0] >> 32)
	valueSize := uint32(result[0])

	// 7- Read the value from the memory
	bytes, ok := module.Memory().Read(valuePosition, valueSize)
	if !ok {
		return nil, errors.New("out of range of memory size")
	}
	return bytes, nil
}
