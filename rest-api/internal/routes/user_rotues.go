package routes

import (
	"net/http"

	"github.com/exzacter/gorestapi/internal/handlers"
)

func SetupUserRoute(mux *http.ServeMux, handler *handlers.Handler) {
	mux.HandleFunc("POST /user/register", handler.CreateUserHandler())
}
