package main

import (
	"net/http"
	"fmt"
	"log"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	statusCode := getStatusCode()
	retryAfter := getRetryAfter(statusCode)
	content := getContentToOutput()
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if isValidStatusCodeForRetryAfter(statusCode) {
			writer.Header().Set("Retry-After", retryAfter)
		}
		writer.WriteHeader(statusCode)
		fmt.Fprint(writer, content)
	})
	log.Fatal(http.ListenAndServe(":80", nil))
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
	statusCodeFromEnvironment := os.Getenv("UNDER_MAINTENANCE_STATUS_CODE")
	if len(statusCodeFromEnvironment) > 0 {
		statusCode, err := strconv.ParseInt(statusCodeFromEnvironment, 10, 64)
		if err != nil {
			log.Printf("'%s' is not a valid status code, defaulting to status 503\n", statusCodeFromEnvironment)
		} else {
			return int(statusCode)
		}
	}
	return http.StatusServiceUnavailable
}

func getRetryAfter(statusCode int) string {
	if isValidStatusCodeForRetryAfter(statusCode) {
		retryAfterFromEnvironment := os.Getenv("UNDER_MAINTENANCE_RETRY_AFTER")
		if len(retryAfterFromEnvironment) > 0 {
			return retryAfterFromEnvironment
		}
	}
	return "300"
}

func isValidStatusCodeForRetryAfter(statusCode int) bool {
	return statusCode == 429 || statusCode == 503
}
