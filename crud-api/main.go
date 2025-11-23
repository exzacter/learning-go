package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// wont be using a DB - using structs and slices. structs are like objects in javascript - all will have key value pairs

// Struct similar to object - for each movie, it has the defining characters as listed. Director is a object of its own as there can be multiple directors. We use a pointer for director the rest are strings. we use `json:""` in order to tag the field and change the way it is displayed within the struct. E.G. goes from ID: 1 to id: 1
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

// similar to above, nothing special
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// creating a variable of movies, that is a slice of the Movie struct. A slice is a dnyamic array, it allows growing and shrinking. We are just initialising it here then appending within func main()
var movies []Movie

// creating a GET function that returns all the movies within the slice as JSON data
func getMovies(w http.ResponseWriter, r *http.Request) {
	// w header adds a header to the server response, within this we are telling the client that the response contains JSON data
	w.Header().Set("Content-Type", "application/json")
	// this turns the slice, in however format it exists in the slice, into a JSON format then sends this to the client
	json.NewEncoder(w).Encode(movies)
}

// returns a single movie by ID from the URL and returns it as JSON
func getMovie(w http.ResponseWriter, r *http.Request) {
	// w header adds a header to the server response, within this we are telling the client that the response contains JSON data
	w.Header().Set("Content-Type", "application/json")
	// Extract the id parameter from the URL (e.g., /movies/1 extracts "1")
	params := mux.Vars(r)
	// _ is a blank identifier, for example here we arent using index
	// loop through all movies, to find the matching ID
	for _, item := range movies {
		if item.ID == params["id"] {
			// convert the matched movie and id to json and send to client
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// create movie reads json input from client, generates a random id for the new movie, adds it to the slice and then returns the new movie
func createMovie(w http.ResponseWriter, r *http.Request) {
	// w header adds a header to the server response, within this we are telling the client that the response contains JSON data
	w.Header().Set("Content-Type", "application/json")
	// create an empty movie struct in order for the client json data to be input
	var movie Movie
	// Read the JSON data from the request body and populate the movie struct
	_ = json.NewDecoder(r.Body).Decode(&movie)
	// Generate a random number and convert it to a string for the movie ID
	movie.ID = strconv.Itoa(rand.Intn(1000000000))
	// add the new movie to the slice
	movies = append(movies, movie)
	// convert the new movie into json and send back to client
	json.NewEncoder(w).Encode(movie)
}

// although incorrect when working in a DB, as we are not in a DB now, we will be deleting the movie and then updating by appending the new movie.
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			// Read the JSON data from the request body and populate the movie struct
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// index is like i = 0 for go
	for index, item := range movies {
		if item.ID == params["id"] {
			// Remove movie at this position and rejoin the slice (deletes the movie)
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func main() {
	// instead of calling mux.NewRouter() each time, just make it r, shortens it.
	r := mux.NewRouter()

	// adding movies to our slice. When we add the movie, we call the slice "movies", then the struct we are following Movie{ID, Isbn, Title, &Director{Firstname, Lastname}}
	// We use the & for director in reference to the pointer. We are telling the append/Movie struct, we are creating the director within the Director Struct, THEN it gets passed to the movie struct, via the pointer as it knows the memory address.
	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie One", Director: &Director{Firstname: "Christopher", Lastname: "Nolan"}})
	movies = append(movies, Movie{ID: "2", Isbn: "438228", Title: "Movie Two", Director: &Director{Firstname: "George", Lastname: "Lucas"}})
	// creating multiple routes via funcs. Specifically different routes allowing different methods
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// printing where the server will be hosted on
	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
