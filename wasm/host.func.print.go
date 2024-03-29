// Package wasm :
package wasm

import (
	"context"
	"fmt"
	"log"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// printStr : print a string to the console
var printStr = api.GoModuleFunc(func(ctx context.Context, module api.Module, params []uint64) {

	// Extract the position and size of the returned value
	// position, length := UnPackPosSize(stack[0])

	position := uint32(params[0])
	length := uint32(params[1])

	buffer, ok := module.Memory().Read(position, length)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", position, length)
	}
	fmt.Println(string(buffer))

	params[0] = 0 // return 0
})

// DefineHostFuncPrintStr defines the hostPrintString function which takes in a
// string position and string length as parameters, and returns an integer.
//
// builder: The HostModuleBuilder object.
//
// Returns: None.
func DefineHostFuncPrintStr(builder wazero.HostModuleBuilder) {
	// hostPrintStr
	builder.NewFunctionBuilder().
		WithGoModuleFunction(printStr,
			[]api.ValueType{
				api.ValueTypeI32, // string position
				api.ValueTypeI32, // string length
			},
			[]api.ValueType{api.ValueTypeI32}).Export("hostPrintStr")

}
