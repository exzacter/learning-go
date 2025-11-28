package main

import (
	"fmt"
	"net/http"
	"log"
)

func main() {

	mux := http.NewServerMux()

	serverAddr := fmt.Sprintf(":%s", config.ServerPort)
	server := &http.Server(
		Addr: serverAddr,
		Handler: nil,
	)

	if err := server.ListenAndServer(); err != nil {
		log.Fatalf("server failed %v", err)
	}
}