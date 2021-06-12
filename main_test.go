package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ohthehugemanatee/zoom-splitter/tools"
)

var cmdLogBuffer tools.TestLogBuffer

func TestRootHandler(t *testing.T) {
	t.Run("Test error response from push URL", func(t *testing.T) {
		logBuffer := tools.CreateAndActivateEmptyTestLogBuffer()
		logBuffer.ExpectLog(LogServerReady)
		logBuffer.ExpectLog(LogRequestReceived)
		logBuffer.ExpectLog(LogNoFileName)
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
		// Without a real file it should return 404.
		AssertStatus(t, http.StatusNotFound, responseRecorder.Code)
	})
	t.Run("Test sending a file to ffmpeg", func(t *testing.T) {
		tmpFile, err := ioutil.TempFile(os.TempDir(), "zoomsplitter-test-")
		if err != nil {
			log.Fatal("Can't create temporary file", err)
		}
		logBuffer := tools.CreateAndActivateEmptyTestLogBuffer()
		logBuffer.ExpectLog(LogServerReady)
		logBuffer.ExpectLog(LogRequestReceived)
		logBuffer.ExpectLog(LogFileRequestReceived + tmpFile.Name())
		logBuffer.ExpectLog(LogFoundFile + tmpFile.Name())
		responseRecorder := runDummyRequest(t, "GET", "/?file="+tmpFile.Name(), RootHandler)
		logBuffer.TestLogValues(t)
		AssertStatus(t, http.StatusOK, responseRecorder.Code)
	})
	t.Run("Test running ffmpeg command", func(t *testing.T) {
		stdout, err := execFfmpeg(fakeExecCommandSuccess)
		if err != nil {
			t.Error(err)
			return
		}
		// Check to make sure the stdout is returned properly
		stdoutStr := stdout.String()
		expected := "echo foobar"
		if stdoutStr != expected {
			t.Errorf("Wrong command executed. Got: \n%s\n Wanted: \n%s", stdoutStr, expected)
		}
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

// This pattern is how exec.Cmd tests itself. Basically create a command context
// which replaces the real exc.Cmd with one which calls a spy function using
// exec.Cmd("go", "-test.run=...")` `. Tricky and a little confusing.

// TestShellProcessSuccess is a method that is called as a substitute for a shell command,
// the GO_TEST_PROCESS flag ensures that if it is not called as part of the test suite, it is
// skipped.
func TestShellProcessSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}
	var cmd []string
	cmd = append(cmd, os.Args[3:]...)
	fmt.Fprint(os.Stdout, strings.Join(cmd, " "))
	os.Exit(0)
}

// fakeExecCommandSuccess is a function that initialises a new exec.Cmd, one which will
// simply call TestShellProcessSuccess rather than the command it is provided. It will
// also pass through the command and its arguments as an argument to TestShellProcessSuccess
func fakeExecCommandSuccess(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestShellProcessSuccess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_TEST_PROCESS=1"}
	return cmd
}
