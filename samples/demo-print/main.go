// Package main :
package main

import (
	"context"
	"log"
	"os"

	"github.com/world-wide-wasm/capsule-host-sdk/wasm"
)

func main() {
	ctx := context.Background()

	wasmFilePath := "../../functions/hello-print/hello-print.wasm"
	wasmFunctionName := "hello"

	runtime, module, errWasm := wasm.Initialize(ctx, wasmFilePath, nil)

	defer runtime.Close(ctx)

	if errWasm != nil {
		os.Exit(1)
	}

	result, err := wasm.CallModuleFunction(module, wasmFunctionName, []byte("Bob Morane"), ctx)

	if err != nil {
		log.Println("😡", err)
	}
	log.Println("🙂", string(result))

}
