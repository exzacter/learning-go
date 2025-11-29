package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/exzacter/gorestapi/internal/auth"
	"github.com/exzacter/gorestapi/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

// creates custom type for context key to avoid collision
type contextKey string

// constant used in storing the user claims
const UserClaimsKey contextKey = "claims"

func AuthMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// retrieves the authorization header from the request (postman/web/mobile)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "No token provided")
			return
		}

		// strips the Bearer string from the bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &auth.Claims{}

		// Parse the token and validating it
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// provided our key from the environment variable and validate it against the token from the request
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		// handle validation error
		if err != nil {
			// handle likely tampered token
			if err == jwt.ErrSignatureInvalid {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid Token")
				return
			}

			// handle any other parsing error - expired, malformed
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid Token")
			return
		}

		// if token is valid, store the claim in the request
		if token.Valid {
			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			r = r.WithContext(ctx) // replace request context with the new request
			next.ServeHTTP(w, r)   // calls the enxt handler, with the updated request
		} else {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid Token")
		}
	})
}
