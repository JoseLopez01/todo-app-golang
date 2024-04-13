package todo

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"todo-app/todo/models"
)

func TestNewRedisRepository(t *testing.T) {
	t.Run("should return a not nil instance", func(t *testing.T) {
		client := redis.NewClient(&redis.Options{})
		repository := NewRedisRepository(client)
		assert.NotNil(t, repository)
		assert.IsType(t, &RedisRepository{}, repository)
	})
}

func TestRedisRepository_Create(t *testing.T) {
	t.Run("should return a nil error and a todo with a uuid", func(t *testing.T) {
		client := getRedisClient(t)
		ctx := context.TODO()
		todo := models.Todo{
			Name:        "name",
			Description: "description",
			StartDate:   time.Now(),
			DueDate:     time.Now().Add(time.Minute * 5),
		}

		repository := NewRedisRepository(client)
		response, err := repository.Create(ctx, "test@test.test", todo)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.ID)
	})

	t.Run("should return an error if the todo cannot be created", func(t *testing.T) {
		client := getRedisClient(t)
		ctx := context.TODO()
		canceled, cancel := context.WithCancel(ctx)
		cancel()
		todo := models.Todo{
			Name:        "name",
			Description: "description",
			StartDate:   time.Now(),
			DueDate:     time.Now().Add(time.Minute * 5),
		}

		repository := NewRedisRepository(client)
		response, err := repository.Create(canceled, "test@test.test", todo)
		assert.ErrorIs(t, err, ErrWhileCreating)
		assert.Zero(t, response)
	})
}

func TestRedisRepository_Delete(t *testing.T) {
	t.Run("should return an error if the id is invalid", func(t *testing.T) {
		client := getRedisClient(t)
		ctx := context.TODO()

		repository := NewRedisRepository(client)
		err := repository.Delete(ctx, "test@test.test", "invalidid")
		assert.ErrorIs(t, err, ErrInvalidID)
	})

	t.Run("should return an error if the todo cannot be deleted", func(t *testing.T) {
		client := getRedisClient(t)
		canceled, cancel := context.WithCancel(context.TODO())
		cancel()

		repository := NewRedisRepository(client)
		err := repository.Delete(canceled, "test@test.test", "279f4a4e-48dc-4569-83df-8b30ce488599")
		assert.ErrorIs(t, err, ErrWhileDeleting)
	})

	t.Run("should return nil if no error happens", func(t *testing.T) {
		client := getRedisClient(t)
		ctx := context.TODO()

		repository := NewRedisRepository(client)
		err := repository.Delete(ctx, "test@test.test", "279f4a4e-48dc-4569-83df-8b30ce488599")
		assert.NoError(t, err)
	})
}

func TestRedisRepository_GetAll(t *testing.T) {
	t.Run("should return an error if the todos can't be retrieved", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		cancel()

		client := getRedisClient(t)
		repository := NewRedisRepository(client)
		response, err := repository.GetAll(ctx, "test@test.test")
		assert.ErrorIs(t, err, ErrWhileRetrieving)
		assert.Nil(t, response)
	})

	t.Run("should return an error if the unmarshal fails", func(t *testing.T) {
		ctx := context.TODO()
		client := getRedisClient(t)
		err := client.HSet(ctx, fmt.Sprintf(redisKey, "test@test.test"), "279f4a4e-48dc-4569-83df-8b30ce488599", "{]").Err()
		assert.NoError(t, err)

		repository := NewRedisRepository(client)
		response, err := repository.GetAll(ctx, "test@test.test")
		assert.ErrorIs(t, err, ErrWhileRetrieving)
		assert.Nil(t, response)
	})

	t.Run("should return the saved todos if no error happens", func(t *testing.T) {
		ctx := context.TODO()
		client := getRedisClient(t)
		todo, err := json.Marshal(models.Todo{
			ID:          "279f4a4e-48dc-4569-83df-8b30ce488599",
			Name:        "name",
			Description: "description",
			StartDate:   time.Now(),
			DueDate:     time.Now().Add(time.Minute * 5),
		})
		assert.NoError(t, err)

		err = client.HSet(ctx, fmt.Sprintf(redisKey, "test@test.test"), "279f4a4e-48dc-4569-83df-8b30ce488599", string(todo)).Err()
		assert.NoError(t, err)

		repository := NewRedisRepository(client)
		response, err := repository.GetAll(ctx, "test@test.test")
		assert.NoError(t, err)
		assert.Len(t, response, 1)
	})
}

func TestRedisRepository_GetByID(t *testing.T) {
	t.Run("should return an error if the id is invalid", func(t *testing.T) {
		client := getRedisClient(t)
		ctx := context.TODO()

		repository := NewRedisRepository(client)
		response, err := repository.GetByID(ctx, "test@test.test", "invalidid")
		assert.ErrorIs(t, err, ErrInvalidID)
		assert.Zero(t, response)
	})

	t.Run("should return an error if the todo can't be retrieved", func(t *testing.T) {
		client := getRedisClient(t)
		canceled, cancel := context.WithCancel(context.TODO())
		cancel()

		repository := NewRedisRepository(client)
		response, err := repository.GetByID(canceled, "test@test.test", "279f4a4e-48dc-4569-83df-8b30ce488599")
		assert.ErrorIs(t, err, ErrWhileRetrieving)
		assert.Zero(t, response)
	})

	t.Run("should return the saved todo", func(t *testing.T) {
		ctx := context.TODO()
		todo, err := json.Marshal(models.Todo{
			ID:          "279f4a4e-48dc-4569-83df-8b30ce488599",
			Name:        "name",
			Description: "description",
			StartDate:   time.Now(),
			DueDate:     time.Now().Add(time.Minute * 5),
		})
		client := getRedisClient(t)
		err = client.HSet(ctx, fmt.Sprintf(redisKey, "test@test.test"), "279f4a4e-48dc-4569-83df-8b30ce488599", string(todo)).Err()
		assert.NoError(t, err)

		repository := NewRedisRepository(client)
		response, err := repository.GetByID(ctx, "test@test.test", "279f4a4e-48dc-4569-83df-8b30ce488599")
		assert.NoError(t, err)
		assert.NotZero(t, response)
		assert.Equal(t, "279f4a4e-48dc-4569-83df-8b30ce488599", response.ID)
	})

	t.Run("should return the unmarshal error", func(t *testing.T) {
		ctx := context.TODO()
		client := getRedisClient(t)
		err := client.HSet(ctx, fmt.Sprintf(redisKey, "test@test.test"), "279f4a4e-48dc-4569-83df-8b30ce488599", `{]`).Err()
		assert.NoError(t, err)

		repository := NewRedisRepository(client)
		response, err := repository.GetByID(ctx, "test@test.test", "279f4a4e-48dc-4569-83df-8b30ce488599")
		assert.ErrorIs(t, err, ErrWhileRetrieving)
		assert.Zero(t, response)
	})
}

func TestRedisRepository_Update(t *testing.T) {
	t.Run("should return an error if the id is invalid", func(t *testing.T) {
		client := getRedisClient(t)
		ctx := context.TODO()

		repository := NewRedisRepository(client)
		response, err := repository.Update(ctx, "test@test.test", "invalidid", models.Todo{})
		assert.ErrorIs(t, err, ErrInvalidID)
		assert.Zero(t, response)
	})

	t.Run("should return an error if the todo can't be updated", func(t *testing.T) {
		client := getRedisClient(t)
		canceled, cancel := context.WithCancel(context.TODO())
		cancel()

		repository := NewRedisRepository(client)
		response, err := repository.Update(canceled, "test@test.test", "279f4a4e-48dc-4569-83df-8b30ce488599", models.Todo{})
		assert.ErrorIs(t, err, ErrWhileUpdating)
		assert.Zero(t, response)
	})

	t.Run("should return the todo if no error happens", func(t *testing.T) {
		client := getRedisClient(t)
		ctx := context.TODO()

		repository := NewRedisRepository(client)
		response, err := repository.Update(ctx, "test@test.test", "279f4a4e-48dc-4569-83df-8b30ce488599", models.Todo{})
		assert.NoError(t, err)
		assert.NotZero(t, response)
		assert.Equal(t, "279f4a4e-48dc-4569-83df-8b30ce488599", response.ID)
	})
}

func getRedisClient(t *testing.T) *redis.Client {
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	ctx := context.TODO()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	redisHost, err := container.Endpoint(ctx, "")
	if err != nil {
		t.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisHost,
	})

	t.Cleanup(func() {
		container.Terminate(ctx)
	})

	return client
}
