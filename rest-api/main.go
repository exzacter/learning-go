package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/exzacter/gorestapi/internal/handlers"
	"github.com/exzacter/gorestapi/internal/routes"
	"github.com/exzacter/gorestapi/serverconfig"
)

func main() {

	// the config I am loading is being imported by the file within serverconfig and function "LoadConfig"
	config, err := serverconfig.LoadConfig()
	// just if there is no config or that is has failed to fetch, print out the error
	if err != nil {
		log.Fatalf("Failed to load config %v", err)
	}

	// thisis calling the core_handler which in future will hold our connections to DB and other things we are dependant on
	handler := handlers.NewHandlers()

	// mux or NewServeMux is the router. It maps the url path from the request and can point them to the function to handle it
	mux := http.NewServeMux()

	// calls the setuproutes function within routes. the setup routes function registers all of the functions and URL's being called within it
	routes.SetupRoutes(mux, handler)

	// setting serverAddr variable to the value of the string from config.ServerPort which is set in the config.go file in serverconfig folder
	serverAddr := fmt.Sprintf(":%s", config.ServerPort)
	// telling server, to run on the port specified in the serverAddr and then all requests to go through mux (router)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}

	fmt.Printf("Server has been started on port %s\n", serverAddr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed %v", err)
	}
}
