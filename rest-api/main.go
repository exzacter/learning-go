package main

import (
	"fmt"
	"net/http"
	"log"

	"github.com/exzacter/gorestapi/serverconfig"
)

func main() {

	config, err := serverconfig.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config %v", err)
	}

	mux := http.NewServerMux()

	mux.Handle("/")

	serverAddr := fmt.Sprintf(":%s", config.ServerPort)
	server := &http.Server(
		Addr: serverAddr,
		Handler: mux,
	)

	if err := server.ListenAndServer(); err != nil {
		log.Fatalf("server failed %v", err)
	}
}
