package routes

import (
	"net/http"

	"github.com/exzacter/gorestapi/internal/handlers"
)

func SetupHealthRoute(mux *http.ServeMux, handler *handlers.Handler) {
	// we are calling our router (mux) to handle the request of "/health" then send it to the correct handler for the right response to be given by the API
	mux.HandleFunc("/health", handler.HealthHandler())
}
