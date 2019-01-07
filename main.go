package main

import (
	"net/http"
	"fmt"
	"log"
	"io/ioutil"
	"os"
	"strconv"
)

var StatusCodeEnvironmentVariableName = "UNDER_MAINTENANCE_STATUS_CODE"
var RetryAfterEnvironmentVariableName = "UNDER_MAINTENANCE_RETRY_AFTER"

var statusCode int
var retryAfter string
var content string

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

func requestHandler(writer http.ResponseWriter, request *http.Request) {
	if isValidStatusCodeForRetryAfter(statusCode) {
		writer.Header().Set("Retry-After", retryAfter)
	}
	writer.WriteHeader(statusCode)
	fmt.Fprint(writer, content)
}

func getContentToOutput() string {
	content := "Under maintenance"
	bytes, err := ioutil.ReadFile("under-maintenance.html")
	if err == nil { // file exists
		log.Println("Found file 'under-maintenance.html', using content of file as output.")
		content = string(bytes)
	} else {
		log.Println("No template file provided, using default output.")
	}
	return content
}

func getStatusCode() int {
	statusCodeFromEnvironment := os.Getenv(StatusCodeEnvironmentVariableName)
	if len(statusCodeFromEnvironment) > 0 {
		statusCode, err := strconv.ParseInt(statusCodeFromEnvironment, 10, 64)
		if err != nil {
			log.Printf("'%s' is not a valid status code, defaulting to %d\n", statusCodeFromEnvironment, http.StatusServiceUnavailable)
		} else {
			return int(statusCode)
		}
	}
	return http.StatusServiceUnavailable
}

func getRetryAfter(statusCode int) string {
	if isValidStatusCodeForRetryAfter(statusCode) {
		retryAfterFromEnvironment := os.Getenv(RetryAfterEnvironmentVariableName)
		if len(retryAfterFromEnvironment) > 0 {
			return retryAfterFromEnvironment
		}
	}
	return "300"
}

func isValidStatusCodeForRetryAfter(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode == http.StatusServiceUnavailable
}
