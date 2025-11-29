package routes

import (
	"net/http"

	"github.com/exzacter/gorestapi/internal/handlers"
	"github.com/exzacter/gorestapi/internal/utils"
)

// passingin the routes that CAN be called via the API so if the client requests them they can be sent to the appropriate handler
func SetupRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	SetupHealthRoute(mux, handler)
	SetupTestRoute(mux, handler)
	SetupUserRoute(mux, handler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		utils.RespondWithNotFound(w)
	})
}
