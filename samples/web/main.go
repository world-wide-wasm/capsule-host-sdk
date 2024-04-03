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
func GetHeader(header map[string][]string) (string, error) {

	// Marshal the map to a JSON byte slice
	jsonData, err := json.Marshal(header)
	if err != nil {
		return fmt.Sprintf(`{"error":"%s"}`, err.Error()), err
	}
	// Convert the byte slice to a string
	jsonString := string(jsonData)

	return jsonString, err
}

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
func GetContentType(header map[string][]string) string {

	value1, ok := header["Content-Type"]
	if ok {
		return value1[0]
	}
	value2, ok := header["content-type"]
	if ok {
		return value2[0]
	}
	return ""
}

/*
ContentType returns a string depending of the content type (text, html and json)
*/
func ContentType(header map[string][]string) string {

	value := strings.Split(GetContentType(header), ";")[0]
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

func main() {

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
		/*
		if request.URL.Path != "/" {
			http.NotFound(response, request)
			return
		}
		*/

		wasm.Protection.Lock()
		// don't forget to release the lock on the Mutex
		defer wasm.Protection.Unlock()

		var body string

		switch request.Method {
		case http.MethodGet:
			body = ""
		case http.MethodPost:
			body = GetBody(request)
		default:
			body = "" // 🚧 work in progress
		}

		header, err := GetHeader(request.Header) // default value of header is "{}" if error
		if err != nil {
			log.Printf("😡 error with the header: %s, default value: %s", err, header)
		}
		contentType := ContentType(request.Header)

		requestData := map[string]string{
			"body":   body,
			"header": header, // -> JSON string
			"type":   contentType,
			"method": request.Method,
			"host": request.Host,
			"remoteAddr": request.RemoteAddr,
			"requestURI": request.RequestURI,
		}

		// create the parameter for the function from the requestData
		funcParam, err := GetJSON(requestData)
		if err != nil {
			log.Printf("😡 error with the function parameter: %s, default value: %s", err, funcParam)
		}

		module := wasm.GetModule()

		// call the wasm function
		result, err := wasm.CallModuleFunction(module, wasmFunctionName, funcParam, ctx)

		if err != nil {
			log.Printf("😡 error with the result of the function: %s", err)
			response.Write([]byte(err.Error()))
		}

		// prepare the response data
		var resultData map[string]interface{}
		err = json.Unmarshal(result, &resultData)
		if err != nil {
			log.Printf("😡 error when unmarshal the result: %s", err)
			response.Write([]byte(err.Error()))
		}

		// read the header from the result of the function
		resultHeader := resultData["header"].(map[string]interface{})
		// set the header od the response with the header from the result of the function
		for key, element := range resultHeader {
			response.Header().Set(key, element.([]interface{})[0].(string))
		}
		// send the response to the client
		response.Write([]byte(resultData["body"].(string)))
	})

	var errListening error
	log.Println("🌍 Capsule http server is listening on: " + httpPort)
	errListening = http.ListenAndServe(":"+httpPort, mux)

	log.Fatal(errListening)

}
