package main

import (
	"fmt"
	"strconv"
	"net/http"
	"html/template"
	"log"
)

// Define a home function which writes a byte slice containing "Hello from Snippetbox" as the response body
func home(w http.ResponseWriter, r *http.Request) {
	// use the header().add() method to add a server: Go header to the response header map. The first paramter is the header name, and the second parameter is the value
	w.Header().Add("Server" ,"Go")

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}
	
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	// extract the value of id wildcard from the request by using r.PathValue() and try and convert it to an integer using the strconv.Atoi() function.
	// If it cant be converted or the value is less than 1, we return a 404 page
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w,r)
		return
	}

	// Use the fmt.Sprintf() function to interpolate the id value with a message, then write it as the HTTP response
	msg := fmt.Sprintf("Display a specific snipper with id %d...", id)
	w.Write([]byte(msg))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet"))
}

// snippetCreatePost handler
func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	
	w.Write([]byte("Save a new snippet..."))
}

