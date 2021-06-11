package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ohthehugemanatee/zoom-splitter/tools"
)

func TestRootHandler(t *testing.T) {
	t.Run("Test error response from push URL", func(t *testing.T) {
		logBuffer := tools.CreateAndActivateEmptyTestLogBuffer()
		logBuffer.ExpectLog(LogServerReady)
		logBuffer.ExpectLog(LogRequestReceived)
		logBuffer.ExpectLog("Request received without a filename. Doing nothing.")
		responseRecorder := runDummyRequest(t, "GET", "/", RootHandler)
		logBuffer.TestLogValues(t)
		AssertStatus(t, http.StatusBadRequest, responseRecorder.Code)
	})
	t.Run("Test read file name value from URL query", func(t *testing.T) {
		filename := TempFileName("test_", "_readFile")
		logBuffer := tools.CreateAndActivateEmptyTestLogBuffer()
		logBuffer.ExpectLog(LogServerReady)
		logBuffer.ExpectLog(LogRequestReceived)
		logBuffer.ExpectLog(LogFileRequestReceived + filename)
		responseRecorder := runDummyRequest(t, "GET", "/?file="+filename, RootHandler)
		logBuffer.TestLogValues(t)
		AssertStatus(t, http.StatusOK, responseRecorder.Code)
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

// TempFileName generates a temporary file name.
func TempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}
