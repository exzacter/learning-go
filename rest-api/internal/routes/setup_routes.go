package routes

import (
	"net/http"

	"github.com/exzacter/gorestapi/internal/handlers"
)

// passingin the routes that CAN be called via the API so if the client requests them they can be sent to the appropriate handler
func SetupRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	SetupHealthRoute(mux, handler)
	SetupTestRoute(mux, handler)
}
