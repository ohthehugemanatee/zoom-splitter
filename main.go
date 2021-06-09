package main

import (
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
	log.Print("Server started, listening on " + networkAddress)
}

// RootHandler handles HTTP requests to /
func RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Received HTTP request to /")
	fileName := r.URL.Query().Get("file")
	if fileName != "" {
		log.Print("File request received: " + fileName)
		w.WriteHeader(http.StatusOK)
		return
	}
	log.Print("Request received without a filename. Doing nothing.")
	w.WriteHeader(http.StatusBadRequest)
}
