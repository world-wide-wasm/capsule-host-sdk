// Package main :
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/world-wide-wasm/capsule-host-sdk/wasm"
)

/*
GetBody returns the body of an HTTP request as a byte slice.

  - It takes a pointer to an http.Request as a parameter.
  - It returns a byte slice.
*/
func GetBody(request *http.Request) []byte {
	body := make([]byte, request.ContentLength)
	request.Body.Read(body)
	return body
}

func main() {

	//fmt.Println(capsule.Version())
	wasmFilePath := "../../functions/hello/hello.wasm"
	//wasmFilePath := "../../functions/hello-print/hello-print.wasm"

	wasmFunctionName := "hello"
	httpPort := "8080"

	mux := http.NewServeMux()

	ctx := context.Background()

	runtime, module, errWasm := wasm.Initialize(ctx, wasmFilePath, nil)
	defer runtime.Close(ctx)
	if errWasm != nil {
		log.Println("😡", errWasm)
		os.Exit(1)
	}
	wasm.StorePlugin(&runtime, &module)

	mux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {

		if request.URL.Path != "/" {
			http.NotFound(response, request)
			return
		}

		switch request.Method {
		case http.MethodGet:
			response.Write([]byte("capsule loves wasm"))

		case http.MethodPost:

			wasm.Protection.Lock()
			// don't forget to release the lock on the Mutex
			defer wasm.Protection.Unlock()

			body := GetBody(request)

			//fmt.Println("✋ Body:", string(body))

			module := wasm.GetModule()

			result, err := wasm.CallModuleFunction(module, wasmFunctionName, body, ctx)

			// TODO
			/*
			- define header
			- content type
			- ...
			*/

			if err != nil {
				response.Write([]byte(err.Error()))
			}
			response.Write(result)

		default:
			//response.Header().Set("Allow", "GET, POST, OPTIONS")
			http.Error(response, "method not allowed", http.StatusMethodNotAllowed)
		}

	})

	var errListening error
	log.Println("🌍 Capsule http server is listening on: " + httpPort)
	errListening = http.ListenAndServe(":"+httpPort, mux)

	log.Fatal(errListening)

}
