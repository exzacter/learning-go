package main

import (
	"fmt"
	"log"
	"net/http"
)

// we will be creating 3 routes for our webserver. "/", "/hello" and "/form"

// within functions that are being created to be served within a HandleFunc we need to provide a writer and reader
// w and r are the variables for ResponseWriter and Request.
// *http.Request is an example of a pointer. This means the "*" is the pointer, and the object is the Request. The pointer allows us to hold the memory of the request address from the http.Request object as opposed to the actual object itself. ELI5 - imagine you have an instruction manual (request URL), someone has requested to see it, so you photocopy each and every page then give a DUPLICATED version of this to the person requesting it. With a pointer, you instead when asked for the book, provide the location (memory location) of the booklet (URL) on a tiny and small notepad, then the person can read the booklet from the original location. TLDR pointer allows us to direct the user to the exact location of the fucntion/object in the memory as opposed to recreating it
func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		// print the error in default or whatever type it uses. %v matches to each argument after the ",". if there is 2 %v's there needs to be 2 outputs after the comma
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful\n")
	// use r for the request of the form that has been posted by the user. then call the name of the field that you want to set
	name := r.FormValue("name")
	address := r.FormValue("address")
	// %s is the way to get the output to be in a string format
	fmt.Fprintf(w, "Name = %s\n", name)
	fmt.Fprintf(w, "Address = %s\n", address)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// r.URL.Path extracts the path from the clientside URL, then the line checks to see if it is NOT hello.
	if r.URL.Path != "/hello" {
		// Error() writes an error response to the client. This uses "w" which is the writer. This would display 404 not found. The StatusNotFound would then provide the status code for the error (404)
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	// now confirming that the Method the user is attempting is GET
	if r.Method != "GET" {
		// if user attempts anything but a GET request, then writer will write "method not supported" and then StatusMethodNotAllowd error status code of (405)
		http.Error(w, "method is not supported", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "hello")
}

func main() {
	// declares and define variables using ":="
	// serving the file using http.FileServer - but specifying the files that are to be served with http.Dir("path/to/file")
	fileServer := http.FileServer(http.Dir("./static"))

	// HandleFunc uses functions that are created, then that is what is served at the handle section in this case "/form" and "/hello"
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)

	// handle "/" route by serving the files that i specified in the FileServer
	// Register the catch-all file server LAST so it doesn't intercept /form and /hello
	http.Handle("/", fileServer)

	fmt.Printf("Starting server at port 8080\n")
	// http.ListenAndServe() the first :8080 or whatever is inside there is the accepting address. As we have stated above we are starting on 8080, we are going to listen on 8080.
	// the secondary option is for a specific handler. if nothing or "nil" then it will use the default handler, this will work because we have registered Handle() and HandleFunc() functions
	// using log.Fatal()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		// fatal prints the error to the console, but then proceeds to kill/terminate the server
		log.Fatal(err)
	}
	fmt.Printf("Server started successfully")
}
