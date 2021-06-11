package main

import (
	"log"
	"net/http"
	"os"
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
	log.Print(LogServerReady)
}

// RootHandler handles HTTP requests to /
func RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Print(LogRequestReceived)
	fileName := r.URL.Query().Get("file")
	if fileName != "" {
		log.Print(LogFileRequestReceived + fileName)
		fileInfo, err := os.Stat("/" + fileName)
		if err == nil {
			log.Print("Found file " + fileInfo.Name())
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Print(LogNoFileName)
	w.WriteHeader(http.StatusBadRequest)
}
