package routes

import (
	"net/http"

	"github.com/exzacter/gorestapi/internal/handlers"
	"github.com/exzacter/gorestapi/internal/middlewares"
)

func SetupUserRoute(mux *http.ServeMux, handler *handlers.Handler) {
	userMux := http.NewServeMux()

	userMux.HandleFunc("POST /register", handler.CreateUserHandler())
	userMux.HandleFunc("POST /login", handler.LoginUserHandler())
	userMux.Handle("GET /profile", middlewares.AuthMiddle(http.HandlerFunc(handler.UserProfile())))

	userMux.Handle("POST /session/logout", middlewares.AuthMiddle(http.HandlerFunc(handler.LogoutHandler())))
	mux.Handle("/users/", http.StripPrefix("/users", userMux))
}
