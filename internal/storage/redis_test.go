package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"todo-app/config"
)

func TestNewRedisClient(t *testing.T) {
	t.Run("should return new redis client", func(t *testing.T) {
		client := NewRedisClient(config.Config{
			RedisHost: "localhost",
			RedisPort: "6379",
		})

		assert.NotNil(t, client)
	})
}
