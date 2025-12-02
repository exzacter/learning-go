package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/exzacter/gorestapi/internal/auth"
	"github.com/exzacter/gorestapi/internal/dtos/request"
	"github.com/exzacter/gorestapi/internal/middlewares"
	"github.com/exzacter/gorestapi/internal/store"
	"github.com/exzacter/gorestapi/internal/utils"
	"github.com/exzacter/gorestapi/internal/validate"
)

// extract function
func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

// clean user session function
func (h *Handler) cleanUserSession(userID string) error {
	// session:123:*
	pattern := fmt.Sprintf("Session:%s:*", userID)

	// background context for redis
	ctx := context.Background()

	// scan to iterate over all the keys matching the pattern declare
	iter := h.Redis.Scan(ctx, 0, pattern, 0).Iterator()

	// loop through each key from redis
	for iter.Next(ctx) {

		// delete the key from redis
		err := h.Redis.Del(ctx, iter.Val()).Err()
		if err != nil {
			fmt.Printf("Failed to delete session")
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}

// logout handler
func (h *Handler) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// extract jwt claims from the context
		claims, ok := r.Context().Value(middlewares.UserClaimsKey).(*auth.Claims)
		if !ok {
			utils.RespondWithError(w, http.StatusBadRequest, "Please login to continue")
			return
		}

		// extract token from auth header
		tokenString := extractTokenFromHeader(r)
		if tokenString == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Missing token")
			return
		}

		// convert expireat to time.Time
		expirationTime := claims.ExpiresAt.Time
		now := time.Now()
		ttl := expirationTime.Sub(now)
		if ttl <= 0 {
			ttl = 5 * time.Minute
		}

		// blacklist token in redis
		err := h.Redis.Set(r.Context(), tokenString, "blacklisted", ttl).Err()
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to blacklist token")
			return
		}

		// clean user session in redis
		userIDStr := fmt.Sprintf("%d", claims.UserID)
		if err := h.cleanUserSession(userIDStr); err != nil {
			fmt.Printf("Error cleaning session for %s: %v", userIDStr, err)
		}

		utils.RespondWithSucess(w, http.StatusOK, "Logged out successfully", true)
	}
}

// profile
func (h *Handler) UserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middlewares.UserClaimsKey).(*auth.Claims)
		if !ok {
			utils.RespondWithError(w, http.StatusBadRequest, "Please login to continue")
			return
		}

		userID := claims.UserID

		// check redis first
		cacheKey := fmt.Sprintf("user:%f", userID)
		if cached, err := h.Redis.Get(r.Context(), cacheKey).Result(); err == nil {
			var user store.User
			if err := json.Unmarshal([]byte(cached), &user); err == nil {
				utils.RespondWithSucess(w, http.StatusOK, "Success (from redis cache)", user)
				return
			}
		}

		// fallback to db
		user, err := h.Queries.GetUser(r.Context(), int32(userID))
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		// set to redis
		userJSON, _ := json.Marshal(user)
		h.Redis.Set(r.Context(), cacheKey, userJSON, 5*time.Minute)

		utils.RespondWithSucess(w, http.StatusOK, "Success", user)
	}
}

// login
func (h *Handler) LoginUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req dtos.LoginRequst
		// user request aka dto
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// validate request
		if err := validate.Validate(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// fetch user from the db using store queries
		user, err := h.Queries.GetUserByUsernameOrEmail(ctx, req.Username)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		if utils.ComparePassword(user.Password, req.Password) {
			utils.RespondWithError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
		token, err := auth.GenerateJWT(int64(user.ID), user.Username, jwtKey)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error generating a token")
			return
		}

		utils.RespondWithSucess(w, http.StatusOK, "Login successful", map[string]string{
			"token": token,
		})

	}
}

// Create user
func (h *Handler) CreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// create context
		ctx := r.Context()

		// user request aka dto
		var req dtos.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// validate request
		if err := validate.Validate(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error while hashing password")
			return
		}
		_, err = h.Queries.CreateUser(ctx, store.CreateUserParams{
			Username: req.Username,
			Email:    req.Email,
			Password: hashedPassword,
		})

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error creating user")
			return
		}

		utils.RespondWithSucess(w, http.StatusCreated, "user created", req.Username)
	}
}
