package todo

import (
	"context"
	"fmt"
	"time"

	"todo-app/todo/dtos"
	"todo-app/todo/models"
)

var (
	ErrInvalidStartDate         = fmt.Errorf("start date can not be parsed")
	ErrInvalidDueDate           = fmt.Errorf("due date can not be parsed")
	ErrStartDateMustBeGTDueDate = fmt.Errorf("start date must be before the due date")
	ErrTodoIsCompleted          = fmt.Errorf("the todo cannot be modified if it's completed")
)

//go:generate mockgen -destination mocks/service_mock.go -package mocks . Service
type Service interface {
	Create(ctx context.Context, email string, dto dtos.CreateTodo) (models.Todo, error)
	GetAll(ctx context.Context, email string) ([]models.Todo, error)
	GetByID(ctx context.Context, email string, id string) (models.Todo, error)
	Delete(ctx context.Context, email string, id string) error
	Update(ctx context.Context, email string, id string, todo dtos.UpdateTodo) (models.Todo, error)
}

type TodosService struct {
	repository Repository
}

func NewTodosService(repository Repository) *TodosService {
	return &TodosService{
		repository: repository,
	}
}

func (t *TodosService) Create(ctx context.Context, email string, dto dtos.CreateTodo) (models.Todo, error) {
	startDate, dueDate, err := validateDates(dto.StartDate, dto.DueDate)
	if err != nil {
		return models.Todo{}, err
	}

	todo := models.Todo{
		DueDate:     dueDate,
		StartDate:   startDate,
		Name:        dto.Name,
		Completed:   false,
		Description: dto.Description,
	}

	return t.repository.Create(ctx, email, todo)
}

func (t *TodosService) GetAll(ctx context.Context, email string) ([]models.Todo, error) {
	return t.repository.GetAll(ctx, email)
}

func (t *TodosService) GetByID(ctx context.Context, email string, id string) (models.Todo, error) {
	return t.repository.GetByID(ctx, email, id)
}

func (t *TodosService) Delete(ctx context.Context, email string, id string) error {
	return t.repository.Delete(ctx, email, id)
}

func (t *TodosService) Update(ctx context.Context, email string, id string, dto dtos.UpdateTodo) (models.Todo, error) {
	startDate, dueDate, err := validateDates(dto.StartDate, dto.DueDate)
	if err != nil {
		return models.Todo{}, err
	}

	todo, err := t.repository.GetByID(ctx, email, id)
	if err != nil {
		return models.Todo{}, err
	}

	if todo.Completed {
		return models.Todo{}, ErrTodoIsCompleted
	}

	todo = models.Todo{
		DueDate:     dueDate,
		StartDate:   startDate,
		Name:        dto.Name,
		Description: dto.Description,
	}

	return t.repository.Update(ctx, email, id, todo)
}

func validateDates(startDate, dueDate string) (time.Time, time.Time, error) {
	parsedStartDate, err := time.Parse(time.DateTime, startDate)
	if err != nil {
		return time.Time{}, time.Time{}, ErrInvalidStartDate
	}

	parsedDueDate, err := time.Parse(time.DateTime, dueDate)
	if err != nil {
		return time.Time{}, time.Time{}, ErrInvalidDueDate
	}

	if parsedStartDate.After(parsedDueDate) {
		return time.Time{}, time.Time{}, ErrStartDateMustBeGTDueDate
	}

	return parsedStartDate, parsedDueDate, nil
}
