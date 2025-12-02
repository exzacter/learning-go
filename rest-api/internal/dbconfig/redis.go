package dbconfig

import (
	"context"
	"fmt"
	"github.com/exzacter/gorestapi/internal/serverconfig"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func ConnectRedis() *redis.Client {
	addr := serverconfig.GetEnv("REDIS_ADDR", "localhost:6379")
	password := serverconfig.GetEnv("REDIS_PASSWORD", "")

	db := 0

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// test connection
	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to redis: %w", err))
	}

	fmt.Println("Connected to redis successfully")

	return rdb
}
