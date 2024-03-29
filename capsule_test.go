// Package capsule
package capsule

import (
	"context"
	"os"
	"testing"

	"github.com/world-wide-wasm/capsule-host-sdk/wasm"
)

func TestGetEnvVarsFromString(t *testing.T) {
	os.Setenv("FIRST_NAME", "Bob")
	os.Setenv("LAST_NAME", "Morane")

	envVarsMap, err := GetEnvVarsFromString(`["FIRST_NAME", "LAST_NAME"]`)

	if (envVarsMap["FIRST_NAME"] == "Bob" && envVarsMap["LAST_NAME"] == "Morane") == false {
		t.Fatalf("Expected: Bob Morane, result: %v %v", envVarsMap["FIRST_NAME"], envVarsMap["LAST_NAME"])
	}

	if err != nil {
		t.Fatalf("Expected: Bob Morane, error: %v", err)
	}

}

func TestGetConfigFromJSONString(t *testing.T) {

	configMap, err := GetConfigFromJSONString(`{"firstName":"Bob", "lastName":"Morane"}`)

	if (configMap["firstName"] == "Bob" && configMap["lastName"] == "Morane") == false {
		t.Fatalf("Expected: Bob Morane, result: %v %v", configMap["firstName"], configMap["lastName"])
	}

	if err != nil {
		t.Fatalf("Expected: Bob Morane, error: %v", err)
	}
}

func TestGetHeaderFromString(t *testing.T) {
	headerName, headerValue := GetHeaderFromString("PRIVATE-TOKEN=1234567890")
	if (headerName == "PRIVATE-TOKEN" && headerValue == "1234567890") == false {
		t.Fatalf("Expected: PRIVATE-TOKEN=1234567890, result: %v %v", headerName, headerValue)
	}
}

// TODO: to be implemented
func TestDownloadWasmFile(t *testing.T) {
	t.Logf("TODO: TestDownloadWasmFile")
}

// 	wasmFilePath := "./functions/hello/hello.wasm"

func TestInitialize(t *testing.T) {
	ctx := context.Background()
	wasmFilePath := "./functions/hello/hello.wasm"

	runtime, module, err := wasm.Initialize(ctx, wasmFilePath, nil)
	if err != nil {
		t.Errorf("Initialize failed: %v", err)
	}

	if runtime == nil {
		t.Errorf("Expected a runtime, got nil")
	}

	if module == nil {
		t.Errorf("Expected a module, got nil")
	}

}

func TestInitialize_ReadFile(t *testing.T) {
	ctx := context.Background()
	wasmFilePath := "nonexistent_file.wasm"

	_, _, err := wasm.Initialize(ctx, wasmFilePath, nil)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}

func TestInitialize_ReadFile_Error(t *testing.T) {
	ctx := context.Background()
	wasmFilePath := "testdata/nonexistent_file.wasm"

	_, _, err := wasm.Initialize(ctx, wasmFilePath, nil)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}


func TestCallModuleFunction(t *testing.T) {
	ctx := context.Background()
	wasmFilePath := "./functions/hello/hello.wasm"
	wasmFunctionName := "hello"


	runtime, module, err := wasm.Initialize(ctx, wasmFilePath, nil)
	if err != nil {
		t.Errorf("Initialize failed: %v", err)
	}

	if runtime == nil {
		t.Errorf("Expected a runtime, got nil")
	}

	if module == nil {
		t.Errorf("Expected a module, got nil")
	}

	result, err := wasm.CallModuleFunction(module, wasmFunctionName, []byte("Bob Morane"), ctx)

	if err != nil {
		t.Errorf("CallModuleFunction failed: %v", err)
	}

	if string(result) != "👋 Hello Bob Morane 😃" {
		t.Fatalf(`Expected: "👋 Hello Bob Morane 😃", result: %v`, string(result))

	}

}

// TODO: test host functions
