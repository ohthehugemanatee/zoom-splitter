package main

import (
	"github.com/ohthehugemanatee/zoom-splitter/logMessage"
	"log"
	"net/http"
)

const networkAddress string = ":80"

// Fatal error handler defined by variable so we can replace it in tests.
var fatalError = log.Fatal

func main() {
	registerHTTPHandlers()
	fatalError(http.ListenAndServe(networkAddress, nil))
}

func registerHTTPHandlers() {
	http.HandleFunc("/", RootHandler)
	log.Print(logMessage.ServerReady)
}

// RootHandler handles HTTP requests to /
func RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Print(logMessage.RequestReceived)
	fileName := r.URL.Query().Get("file")
	if fileName != "" {
		log.Print(logMessage.FileRequestReceived + fileName)
		w.WriteHeader(http.StatusOK)
		return
	}
	log.Print("Request received without a filename. Doing nothing.")
	w.WriteHeader(http.StatusBadRequest)
}
