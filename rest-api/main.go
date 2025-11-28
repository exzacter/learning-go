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

	config, err := serverconfig.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config %v", err)
	}

	handler := handlers.NewHandlers()

	mux := http.NewServerMux()

	routes.SetupRoutes(mux, handler)

	serverAddr := fmt.Sprintf(":%s", config.ServerPort)
	server := &http.Server(
		Addr: serverAddr,
		Handler: mux,
	)

	if err := server.ListenAndServer(); err != nil {
		log.Fatalf("server failed %v", err)
	}
}
