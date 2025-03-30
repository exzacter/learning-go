package main

import (
	"log"
	"net/http"
)

func main() {
	// Use the http.NewServeMux() function to initialise a new servemu, then register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home) // restrict this route to exzact matches on / only
	mux.HandleFunc("GET /snippet/view/{id}", snippetView) // add the {id} wildcard segment
	mux.HandleFunc("GET /snippet/create", snippetCreate)

	// create the new route which is restricted to POST request only
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// Print a log message to say that the server is starting
	log.Print("Starting server on :8008")

	// Use the ListenAndServe() function to start the web server. We pass in two parameters: TCP port to listen on, and the servemux we just created.
	// we use the Log.Fatal() to log the error message and exit.
	err := http.ListenAndServe(":8008", mux)
	log.Fatal(err)
}
