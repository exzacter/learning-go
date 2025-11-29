package handlers

import (
	"database/sql"

	"github.com/exzacter/gorestapi/internal/store"
)

type Handler struct {
	// DB instanc
	DB *sql.DB
	// Query stores
	Queries store.Queries
}

func NewHandlers(db *sql.DB, queries *store.Queries) *Handler {
	return &Handler{
		DB:      db,
		Queries: *queries,
	}
}
