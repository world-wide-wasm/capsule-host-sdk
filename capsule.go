// Package capsule : this is the main Package
package capsule

import (
	_ "embed"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

//go:embed version.txt
var version []byte

// Version returns the version number of the Capsule Host SDK
func Version() string {
	return string(version)
}

// GetEnvVarsFromString :
// returns a map of the forwarded environment variables
func GetEnvVarsFromString(envars string) (map[string]string, error) {
	var vars []string
	envVarsMap := make(map[string]string)

	unmarshallError := json.Unmarshal([]byte(envars), &vars)
	if unmarshallError != nil {
		return nil, unmarshallError
	}
	if len(vars) > 0 {
		for _, envVar := range vars {
			envVarsMap[envVar] = os.Getenv(envVar)
		}
	}

	return envVarsMap, nil
}

// GetConfigFromJSONString :
// read the JSON config string
func GetConfigFromJSONString(config string) (map[string]string, error) {
	var configMap map[string]string
	unmarshallError := json.Unmarshal([]byte(config), &configMap)
	/*
	if unmarshallError != nil {
		fmt.Println("😡 getConfigFromJSONString:", unmarshallError)
		os.Exit(1)
	}
	*/
	return configMap, unmarshallError
}


// GetHeaderFromString :
func GetHeaderFromString(headerNameAndValue string) (string, string) {
	splitHeader := strings.Split(headerNameAndValue, "=")
	headerName := splitHeader[0]
	// join all item of splitAuthHeader with "" except the first one
	headerValue := strings.Join(splitHeader[1:], "")
	return headerName, headerValue
}

// DownloadWasmFile :
// download a wasm remote file
func DownloadWasmFile(url, urlAuthHeader, wasmFilePath string) error {
	// authenticationHeader:
	// Example: "PRIVATE-TOKEN=${GITLAB_WASM_TOKEN}"
	client := resty.New()

	if urlAuthHeader != "" {
		authHeaderName, authHeaderValue := GetHeaderFromString(urlAuthHeader)
		client.SetHeader(authHeaderName, authHeaderValue)

	}

	resp, err := client.R().
		SetOutput(wasmFilePath).
		Get(url)

	if resp.IsError() {
		return errors.New("Error while downloading the wasm file, you should check the authentication token")
	}

	if err != nil {
		return err
	}
	return nil
}
