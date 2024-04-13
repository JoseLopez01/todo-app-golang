package todo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"todo-app/todo/models"
)

const (
	redisKey = "todo-%s"
)

var (
	ErrWhileCreating   = fmt.Errorf("error while creating")
	ErrWhileRetrieving = fmt.Errorf("error while retreving")
	ErrWhileDeleting   = fmt.Errorf("error while deleting")
	ErrWhileUpdating   = fmt.Errorf("error while updating")
	ErrInvalidID       = fmt.Errorf("invalid id")
)

//go:generate mockgen -destination mocks/repository_mock.go -package mocks . Repository
type Repository interface {
	Create(ctx context.Context, email string, todo models.Todo) (models.Todo, error)
	GetAll(ctx context.Context, email string) ([]models.Todo, error)
	GetByID(ctx context.Context, email string, id string) (models.Todo, error)
	Delete(ctx context.Context, email string, id string) error
	Update(ctx context.Context, email string, id string, todo models.Todo) (models.Todo, error)
}

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

func (r *RedisRepository) Create(ctx context.Context, email string, todo models.Todo) (models.Todo, error) {
	todo.ID = uuid.NewString()
	todoBytes, err := json.Marshal(todo)
	if err != nil {
		return models.Todo{}, ErrWhileCreating
	}

	userKey := fmt.Sprintf(redisKey, email)
	err = r.client.HSet(ctx, userKey, todo.ID, todoBytes).Err()
	if err != nil {
		return models.Todo{}, ErrWhileCreating
	}

	return todo, nil
}

func (r *RedisRepository) GetAll(ctx context.Context, email string) ([]models.Todo, error) {
	userKey := fmt.Sprintf(redisKey, email)
	result, err := r.client.HGetAll(ctx, userKey).Result()
	if err != nil {
		return nil, ErrWhileRetrieving
	}

	var todos []models.Todo
	for _, todoString := range result {
		var todo models.Todo
		if err = json.Unmarshal([]byte(todoString), &todo); err != nil {
			return nil, ErrWhileRetrieving
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *RedisRepository) GetByID(ctx context.Context, email string, id string) (models.Todo, error) {
	if err := validateID(id); err != nil {
		return models.Todo{}, err
	}

	userKey := fmt.Sprintf(redisKey, email)
	result, err := r.client.HGet(ctx, userKey, id).Result()
	if err != nil {
		return models.Todo{}, ErrWhileRetrieving
	}

	var todo models.Todo
	if err = json.Unmarshal([]byte(result), &todo); err != nil {
		return models.Todo{}, ErrWhileRetrieving
	}

	return todo, nil
}

func (r *RedisRepository) Delete(ctx context.Context, email string, id string) error {
	if err := validateID(id); err != nil {
		return err
	}

	userKey := fmt.Sprintf(redisKey, email)
	err := r.client.HDel(ctx, userKey, id).Err()
	if err != nil {
		return ErrWhileDeleting
	}

	return nil
}

func (r *RedisRepository) Update(ctx context.Context, email string, id string, todo models.Todo) (models.Todo, error) {
	if err := validateID(id); err != nil {
		return models.Todo{}, err
	}

	todoBytes, err := json.Marshal(todo)
	if err != nil {
		return models.Todo{}, ErrWhileUpdating
	}

	userKey := fmt.Sprintf(redisKey, email)
	err = r.client.HSet(ctx, userKey, id, todoBytes).Err()
	if err != nil {
		return models.Todo{}, ErrWhileUpdating
	}

	todo.ID = id
	return todo, err
}

func validateID(id string) error {
	if err := uuid.Validate(id); err != nil {
		return ErrInvalidID
	}

	return nil
}
