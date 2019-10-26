package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func resetTestEnvironment() {
	_ = os.Unsetenv(StatusCodeEnvironmentVariableName)
	_ = os.Unsetenv(RetryAfterEnvironmentVariableName)
	_ = os.Remove("under-maintenance.html")
}

func TestDefaultStatusCode(t *testing.T) {
	resetTestEnvironment()
	LoadConfiguration()

	w, resp := createTestServer(t, "/")

	if resp.StatusCode != 503 {
		t.Errorf("Status code should be 503 by default, but was %d instead", resp.StatusCode)
	}
	if resp.Header.Get("Retry-After") != "300" {
		t.Errorf("Retry-After header value should be 300 by default, but was %s instead", resp.Header.Get("Retry-After"))
	}
	w.Flush()
}

func TestCustomStatusCode(t *testing.T) {
	resetTestEnvironment()
	_ = os.Setenv(StatusCodeEnvironmentVariableName, "200")
	LoadConfiguration()

	w, resp := createTestServer(t, "/")

	if resp.StatusCode != 200 {
		t.Errorf("Status code should be 200, but was %d instead", resp.StatusCode)
	}
	if resp.Header.Get("Retry-After") != "" && isValidStatusCodeForRetryAfter(resp.StatusCode) {
		t.Errorf("There should be no Retry-After header, because %d is not a valid status code for said header", resp.StatusCode)
	}
	w.Flush()
}

func TestDefaultMaintenancePage(t *testing.T) {
	resetTestEnvironment()
	LoadConfiguration()

	w, resp := createTestServer(t, "/")

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "Under maintenance" {
		t.Errorf("Response body does not match: got %s, expected %s", string(output), "Under maintenance")
	}
	w.Flush()
}

func TestCustomMaintenancePage(t *testing.T) {
	resetTestEnvironment()
	f, err := os.Create("under-maintenance.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, _ = f.WriteString("test")
	LoadConfiguration()

	w, resp := createTestServer(t, "/")

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "test" {
		t.Errorf("Response body does not match: got %s, expected %s", string(output), "test")
	}
	w.Flush()
}

func TestNonRootPathIsHandled(t *testing.T) {
	resetTestEnvironment()
	LoadConfiguration()

	w, resp := createTestServer(t, "/some/other/path")

	if resp.StatusCode != 503 {
		t.Errorf("Status code should be 503 by default, but was %d instead", resp.StatusCode)
	}
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "Under maintenance" {
		t.Errorf("Response body does not match: got %s, expected %s", string(output), "Under maintenance")
	}
	w.Flush()
}

////////////////////////
// NON-TEST FUNCTIONS //
////////////////////////

func createTestServer(t *testing.T, path string) (*httptest.ResponseRecorder, *http.Response) {
	r, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(requestHandler)
	handler.ServeHTTP(w, r)
	resp := w.Result()
	return w, resp
}
