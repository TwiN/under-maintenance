package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	StatusCodeEnvironmentVariableName = "UNDER_MAINTENANCE_STATUS_CODE"
	RetryAfterEnvironmentVariableName = "UNDER_MAINTENANCE_RETRY_AFTER"
	CustomFileEnvironmentVariableName = "UNDER_MAINTENANCE_CUSTOM_FILE_PATH"
)

var (
	statusCode int
	retryAfter string
	content    string
)

func init() {
	// Because the config needs to be reloaded for every tests, all initialized
	// variables are moved to a different function instead of init()
	LoadConfiguration()
}

func main() {
	http.HandleFunc("/", requestHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func LoadConfiguration() {
	statusCode = getStatusCode()
	retryAfter = getRetryAfter(statusCode)
	content = getContentToOutput()
}

func requestHandler(writer http.ResponseWriter, _ *http.Request) {
	if isValidStatusCodeForRetryAfter(statusCode) {
		writer.Header().Set("Retry-After", retryAfter)
	}
	writer.WriteHeader(statusCode)
	fmt.Fprint(writer, content)
}

func getContentToOutput() string {
	content := "Under maintenance"
	filePath := "under-maintenance.html"
	if customFilePath := os.Getenv(CustomFileEnvironmentVariableName); len(customFilePath) > 0 {
		filePath = customFilePath
	}
	if bytes, err := ioutil.ReadFile(filePath); err == nil { // file exists
		content = string(bytes)
		log.Printf("Found file '%s', using content of file as output:\n%s", filePath, string(bytes))
	} else {
		log.Println("No template file provided, using default output.")
	}
	return content
}

func getStatusCode() int {
	if statusCodeFromEnvironment := os.Getenv(StatusCodeEnvironmentVariableName); len(statusCodeFromEnvironment) > 0 {
		if statusCode, err := strconv.ParseInt(statusCodeFromEnvironment, 10, 64); err != nil {
			log.Printf("'%s' is not a valid status code, defaulting to %d\n", statusCodeFromEnvironment, http.StatusServiceUnavailable)
		} else {
			return int(statusCode)
		}
	}
	return http.StatusServiceUnavailable
}

func getRetryAfter(statusCode int) string {
	if isValidStatusCodeForRetryAfter(statusCode) {
		if retryAfterFromEnvironment := os.Getenv(RetryAfterEnvironmentVariableName); len(retryAfterFromEnvironment) > 0 {
			return retryAfterFromEnvironment
		}
	}
	return "300"
}

func isValidStatusCodeForRetryAfter(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode == http.StatusServiceUnavailable
}
