package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/exzacter/gorestapi/internal/auth"
	"github.com/exzacter/gorestapi/internal/dtos/request"
	"github.com/exzacter/gorestapi/internal/middlewares"
	"github.com/exzacter/gorestapi/internal/store"
	"github.com/exzacter/gorestapi/internal/utils"
	"github.com/exzacter/gorestapi/internal/validate"
)

// profile
func (h *Handler) UserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middlewares.UserClaimsKey).(*auth.Claims)
		if !ok {
			utils.RespondWithError(w, http.StatusBadRequest, "Please login to continue")
			return
		}

		userID := claims.UserID

		user, err := h.Queries.GetUser(r.Context(), int32(userID))
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

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
