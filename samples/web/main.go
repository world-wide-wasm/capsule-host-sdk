// Package main :
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/world-wide-wasm/capsule-host-sdk/wasm"
)

/*
GetBody returns the body of an HTTP request as a string.

  - It takes a pointer to an http.Request as a parameter.
  - It returns a string.
*/
func GetBody(request *http.Request) string { // what todo if it's a json string?
	body := make([]byte, request.ContentLength)
	request.Body.Read(body)
	return string(body)
}

/*
GetHeader returns the header od an HTTP request
*/
func GetHeader(request *http.Request) (string, error) {

	// Marshal the map to a JSON byte slice
	jsonData, err := json.Marshal(request.Header)
	if err != nil {
		return fmt.Sprintf(`{"error":"%s"}`, err.Error()), err
	}
	// Convert the byte slice to a string
	jsonString := string(jsonData)

	return jsonString, err
}

// datetime := fmt.Sprintf("%s,%s", date, time)

/*
GetJSON returns a JSON byte slice form a map
*/
func GetJSON(data map[string]string) ([]byte, error) {
	// Marshal the map to a JSON byte slice
	jsonData, err := json.Marshal(data)
	if err != nil {
		return []byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())), err
	}
	return jsonData, nil
}

/*
GetJSONString returns a JSON string form a map
*/
func GetJSONString(data map[string]string) (string, error) {
	jsonData, err := GetJSON(data)
	if err != nil {
		return string(jsonData), err
	}
	// Convert the byte slice to a string
	jsonString := string(jsonData)
	return jsonString, nil
}

/*
GetContentType returns a string with the content type from the request
*/
func GetContentType(request *http.Request) string {

	value1, ok := request.Header["Content-Type"]
	if ok {
		return value1[0]
	}
	value2, ok := request.Header["content-type"]
	if ok {
		return value2[0]
	}
	return ""
}

/*
ContentType returns a string depending of the content type (text, html and json)
*/
func ContentType(request *http.Request) string {

	value := strings.Split(GetContentType(request), ";")[0]
	//fmt.Println("✋", value)
	switch {
	case value == "text/plain":
		return "text"
	case value == "text/html":
		return "html"
	case value == "application/json":
		return "json"
	default:
		return value
	}
}
/*
text/plain: Plain text data (e.g., log messages).
text/html: HTML code for web pages.
application/json: JSON formatted data (commonly used for APIs).
*/

func main() {

	//fmt.Println(capsule.Version())
	//wasmFilePath := "../../functions/hello/hello.wasm"
	//wasmFilePath := "../../functions/hello-print/hello-print.wasm"
	wasmFilePath := "../../functions/nano-service/nano-service.wasm"

	wasmFunctionName := "functionHandler"
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
			header, err := GetHeader(request) // default value of header is "{}" if error
			if err != nil {
				log.Printf("😡 error with the header: %s, default value: %s", err, header)
			}
			contentType := ContentType(request)

			// ✋ depending the content type I need to encode the body string (text/html/json ...)
			// 🚧 TODO: method to get the content type

			/*
				Header = map[string][]string{
					"Accept-Encoding": {"gzip, deflate"},
					"Accept-Language": {"en-us"},
					"Foo": {"Bar", "two"},
				}
			*/

			requestData := map[string]string{
				"body":   body,
				"header": header, // -> JSON string
				"type": contentType,
			}

			funcParam, err := GetJSON(requestData)
			if err != nil {
				log.Printf("😡 error with the function parameter: %s, default value: %s", err, funcParam)
			}

			//fmt.Println("✋ Body:", string(body))

			module := wasm.GetModule()

			result, err := wasm.CallModuleFunction(module, wasmFunctionName, funcParam, ctx)

			// TODO
			/*
				- define header
				- content type -> GetContentType
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
