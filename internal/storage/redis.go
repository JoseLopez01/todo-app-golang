package storage

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	"todo-app/config"
)

func NewRedisClient(config config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: "",
		DB:       0,
	})

	return client
}
