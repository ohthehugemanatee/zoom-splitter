package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ohthehugemanatee/zoom-splitter/tools"
)

func TestRootHandler(t *testing.T) {
	t.Run("Test error response from push URL", func(t *testing.T) {
		logBuffer := tools.CreateAndActivateEmptyTestLogBuffer()
		logBuffer.ExpectLog("Received HTTP request to /")
		responseRecorder := runDummyRequest(t, "GET", "/", RootHandler)
		logBuffer.TestLogValues(t)
		AssertStatus(t, http.StatusNotImplemented, responseRecorder.Code)
	})
}

func runDummyRequest(t *testing.T, verb string, path string, handlerFunc func(w http.ResponseWriter, r *http.Request)) httptest.ResponseRecorder {
	http.DefaultServeMux = http.NewServeMux()
	registerHTTPHandlers()
	request, err := http.NewRequest(verb, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	responseRecorder := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(responseRecorder, request)
	return *responseRecorder
}

// AssertStatus is a test convenience function to compare HTTP status codes.
func AssertStatus(t *testing.T, expected int, got int) {
	if got != expected {
		t.Errorf("Got wrong status code: got %v want %v",
			got, expected)
	}
}
