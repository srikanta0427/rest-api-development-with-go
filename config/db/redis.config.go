package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
func CheckRedis(client *redis.Client) error {
	return client.Ping(context.Background()).Err()
}
