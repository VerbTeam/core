package database

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func Init() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_PUBLIC_ENDPOINT"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORDS"),
		DB:       0,
	})
}
