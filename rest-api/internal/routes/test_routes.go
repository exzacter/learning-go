package routes

import (
	"net/http"

	"github.com/exzacter/gorestapi/internal/handlers"
)

func SetupTestRoute(mux *http.ServeMux, handler *handlers.Handler) {
	// we are calling our router (mux) to handle the request of "/test" then send it to the correct handler for the right response to be given by the API
	mux.HandleFunc("/test", handler.TestHandler())
}
