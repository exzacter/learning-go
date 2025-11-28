package routes

import (
	"net/http"

	"github.com/exzacter/gorestapi/internal/handlers"
)

func SetupRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	SetupHealthRoute(mux, handler)
}
