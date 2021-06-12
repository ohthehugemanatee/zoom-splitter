package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type execContext = func(name string, arg ...string) *exec.Cmd

const networkAddress string = ":80"

// Fatal error handler defined by variable so we can replace it in tests.
var fatalError = log.Fatal

func main() {
	registerHTTPHandlers()
	execFfmpeg(exec.Command)
	fatalError(http.ListenAndServe(networkAddress, nil))
}

func registerHTTPHandlers() {
	http.HandleFunc("/", RootHandler)
	log.Print(LogServerReady)
}

// RootHandler handles HTTP requests to /
func RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Print(LogRequestReceived)
	fileName := r.URL.Query().Get("file")
	if fileName != "" {
		log.Print(LogFileRequestReceived + fileName)
		_, err := os.Stat("/" + fileName)
		if err == nil {
			log.Print(LogFoundFile + fileName)
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Print(LogNoFileName)
	w.WriteHeader(http.StatusBadRequest)
}

func execFfmpeg(cmdContext execContext) (*bytes.Buffer, error) {
	cmd := cmdContext("echo", "foobar")

	// Set up byte buffers to read stdout
	var output bytes.Buffer
	cmd.Stdout = &output

	err := cmd.Run()
	if err != nil {
		return &output, err
	}

	return &output, nil
}
