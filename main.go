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
	log.Print("Server started, listening on " + networkAddress)
}

func registerHTTPHandlers() {
	http.HandleFunc("/", RootHandler)
}

// RootHandler handles HTTP requests to /
func RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Received HTTP request to /")
	w.WriteHeader(http.StatusNotImplemented)
}
