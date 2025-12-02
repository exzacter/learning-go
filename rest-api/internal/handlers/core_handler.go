package handlers

import (
	"database/sql"

	"github.com/exzacter/gorestapi/internal/store"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	// DB instanc
	DB *sql.DB
	// Query stores
	Queries *store.Queries
	Redis   *redis.Client
}

func NewHandlers(db *sql.DB, queries *store.Queries, redisClient *redis.Client) *Handler {
	return &Handler{
		DB:      db,
		Queries: queries,
		Redis:   redisClient,
	}
}
