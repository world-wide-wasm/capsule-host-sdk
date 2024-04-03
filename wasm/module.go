// Package wasm :
package wasm

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// Plugin : specific type
type Plugin struct {
	WazeroModule  api.Module
	WazeroRuntime wazero.Runtime
}

// ModuleFunction : used to add host functions
type ModuleFunction struct {
	GoModuleFunc api.GoModuleFunc
	Params []byte
	Results []byte
}

var (
	plugins = make(map[string]*Plugin)
	// Protection is a mutex to protect the module
	Protection sync.Mutex
)

/*
StorePlugin stores the Wazero runtime and the wasm module into maps
*/
func StorePlugin(runtime *wazero.Runtime, module *api.Module) {
	plugins["main"] = &Plugin{
		WazeroModule:  *module,
		WazeroRuntime: *runtime,
	}
}

/*
GetModule returns the wasm module
*/
func GetModule() api.Module {
	return *&plugins["main"].WazeroModule
}

/*
Initialize initializes the Wasm Plugin
*/
func Initialize(ctx context.Context, wasmFilePath string, goModFuncs []ModuleFunction) (wazero.Runtime, api.Module, error) {
	// 1- Create instance of a wazero runtime
	// Create a new runtime.
	runtime := wazero.NewRuntime(ctx)

	// This closes everything this Runtime created.
	//defer runtime.Close(ctx)

	// The builder allows to add host functions
	// Instantiate a Go-defined module named "env"

	builder := runtime.NewHostModuleBuilder("env")

	// Default host functions
	DefineHostFuncPrintStr(builder)

	if goModFuncs != nil {
		// TODO: add user defined host functions
	}

	// Instantiate builder and default host functions
	_, err := builder.Instantiate(ctx)
	if err != nil {
		log.Panicln("Error with env module and host function(s):", err)
	}

	// Instantiate WASI
	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	// 2- Load the WebAssembly module
	wasmFile, errReadFile := os.ReadFile(wasmFilePath)
	if errReadFile != nil {
		return nil, nil, errReadFile
	}

	// 3- Instantiate the Wasm plugin/program.
	module, errModule := runtime.Instantiate(ctx, wasmFile)
	// temp
	if errModule != nil {
		return nil, nil, errModule
	}

	return runtime, module, nil
}
