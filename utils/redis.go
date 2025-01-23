package utils

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var (
	Ctx    = context.Background()
	Client = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
)
