package todo

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"todo-app/todo/dtos"
	"todo-app/todo/mocks"
	"todo-app/todo/models"
)

func TestNewTodosService(t *testing.T) {
	t.Run("should return a not nil instance", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repository := mocks.NewMockRepository(ctrl)
		service := NewTodosService(repository)
		assert.NotNil(t, service)
		assert.IsType(t, &TodosService{}, service)
	})
}

func TestTodosService_Create(t *testing.T) {
	email := "test@test.test"
	ctx := context.TODO()
	ctxMatcher := reflect.TypeOf((*context.Context)(nil)).Elem()
	validStartDate := time.Now().Format(time.DateTime)
	validDueDate := time.Now().Add(time.Minute * 5).Format(time.DateTime)

	t.Run("should return the validateDates error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repository := mocks.NewMockRepository(ctrl)
		dto := dtos.CreateTodo{
			DueDate:     "invalid",
			StartDate:   "invalid",
			Description: "description",
			Name:        "name",
		}

		service := NewTodosService(repository)
		response, err := service.Create(ctx, email, dto)
		assert.Error(t, err)
		assert.Zero(t, response)
	})

	t.Run("should return the repository response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repository := mocks.NewMockRepository(ctrl)
		repository.
			EXPECT().
			Create(gomock.AssignableToTypeOf(ctxMatcher), gomock.Eq(email), gomock.AssignableToTypeOf(models.Todo{})).
			Return(models.Todo{
				ID: "279f4a4e-48dc-4569-83df-8b30ce488599",
			}, nil)

		dto := dtos.CreateTodo{
			DueDate:     validDueDate,
			StartDate:   validStartDate,
			Description: "description",
			Name:        "name",
		}

		service := NewTodosService(repository)
		response, err := service.Create(ctx, email, dto)
		assert.NoError(t, err)
		assert.NotZero(t, response)
	})
}

func TestTodosService_Delete(t *testing.T) {
	ctx := context.TODO()
	ctxMatcher := reflect.TypeOf((*context.Context)(nil)).Elem()
	email := "test@test.test"
	id := "279f4a4e-48dc-4569-83df-8b30ce488599"

	t.Run("should delete the todo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repository := mocks.NewMockRepository(ctrl)
		repository.
			EXPECT().
			Delete(gomock.AssignableToTypeOf(ctxMatcher), gomock.Eq(email), gomock.Eq(id)).
			Return(nil)

		service := NewTodosService(repository)
		err := service.Delete(ctx, email, id)
		assert.NoError(t, err)
	})
}

func TestTodosService_GetAll(t *testing.T) {
	ctx := context.TODO()
	ctxMatcher := reflect.TypeOf((*context.Context)(nil)).Elem()
	email := "test@test.test"
	id := "279f4a4e-48dc-4569-83df-8b30ce488599"

	t.Run("should get all the todos", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repository := mocks.NewMockRepository(ctrl)
		repository.
			EXPECT().
			GetAll(gomock.AssignableToTypeOf(ctxMatcher), gomock.Eq(email)).
			Return([]models.Todo{{ID: id}}, nil)

		service := NewTodosService(repository)
		response, err := service.GetAll(ctx, email)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
	})
}

func TestTodosService_GetByID(t *testing.T) {
	ctx := context.TODO()
	ctxMatcher := reflect.TypeOf((*context.Context)(nil)).Elem()
	email := "test@test.test"
	id := "279f4a4e-48dc-4569-83df-8b30ce488599"

	t.Run("should get by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repository := mocks.NewMockRepository(ctrl)
		repository.
			EXPECT().
			GetByID(gomock.AssignableToTypeOf(ctxMatcher), gomock.Eq(email), gomock.Eq(id)).
			Return(models.Todo{ID: id}, nil)

		service := NewTodosService(repository)
		response, err := service.GetByID(ctx, email, id)
		assert.NoError(t, err)
		assert.Equal(t, id, response.ID)
	})
}

func TestTodosService_Update(t *testing.T) {
	email := "test@test.test"
	ctx := context.TODO()
	ctxMatcher := reflect.TypeOf((*context.Context)(nil)).Elem()
	validStartDate := time.Now().Format(time.DateTime)
	validDueDate := time.Now().Add(time.Minute * 5).Format(time.DateTime)
	id := "279f4a4e-48dc-4569-83df-8b30ce488599"

	t.Run("should return the validateDates error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repository := mocks.NewMockRepository(ctrl)
		dto := dtos.UpdateTodo{
			DueDate:     "invalid",
			StartDate:   "invalid",
			Description: "description",
			Name:        "name",
		}

		service := NewTodosService(repository)
		response, err := service.Update(ctx, email, id, dto)
		assert.Error(t, err)
		assert.Zero(t, response)
	})

	t.Run("should return the GetById error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repository := mocks.NewMockRepository(ctrl)
		repository.
			EXPECT().
			GetByID(gomock.AssignableToTypeOf(ctxMatcher), gomock.Eq(email), gomock.Eq(id)).
			Return(models.Todo{}, fmt.Errorf("error"))

		dto := dtos.UpdateTodo{
			DueDate:     validDueDate,
			StartDate:   validStartDate,
			Description: "description",
			Name:        "name",
		}

		service := NewTodosService(repository)
		response, err := service.Update(ctx, email, id, dto)
		assert.Error(t, err)
		assert.Zero(t, response)
	})

	t.Run("should return the ErrTodoIsCompleted error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repository := mocks.NewMockRepository(ctrl)
		repository.
			EXPECT().
			GetByID(gomock.AssignableToTypeOf(ctxMatcher), gomock.Eq(email), gomock.Eq(id)).
			Return(models.Todo{Completed: true}, nil)

		dto := dtos.UpdateTodo{
			DueDate:     validDueDate,
			StartDate:   validStartDate,
			Description: "description",
			Name:        "name",
		}

		service := NewTodosService(repository)
		response, err := service.Update(ctx, email, id, dto)
		assert.ErrorIs(t, err, ErrTodoIsCompleted)
		assert.Zero(t, response)
	})

	t.Run("should update the todo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repository := mocks.NewMockRepository(ctrl)
		repository.
			EXPECT().
			GetByID(gomock.AssignableToTypeOf(ctxMatcher), gomock.Eq(email), gomock.Eq(id)).
			Return(models.Todo{Completed: false, ID: id}, nil)
		repository.
			EXPECT().
			Update(gomock.AssignableToTypeOf(ctxMatcher), gomock.Eq(email), gomock.Eq(id), gomock.AssignableToTypeOf(models.Todo{})).
			Return(models.Todo{ID: id}, nil)

		dto := dtos.UpdateTodo{
			DueDate:     validDueDate,
			StartDate:   validStartDate,
			Description: "description",
			Name:        "name",
		}

		service := NewTodosService(repository)
		response, err := service.Update(ctx, email, id, dto)
		assert.NoError(t, err, ErrTodoIsCompleted)
		assert.NotZero(t, response)
		assert.Equal(t, id, response.ID)
	})
}

func Test_ValidateDates(t *testing.T) {
	validStartDate := time.Now().Format(time.DateTime)
	validDueDate := time.Now().Add(time.Minute * 5).Format(time.DateTime)

	t.Run("should return the ErrInvalidStartDate", func(t *testing.T) {
		_, _, err := validateDates("invalid", validDueDate)
		assert.ErrorIs(t, err, ErrInvalidStartDate)
	})

	t.Run("should return the ErrInvalidDueDate", func(t *testing.T) {
		_, _, err := validateDates(validStartDate, "invalid")
		assert.ErrorIs(t, err, ErrInvalidDueDate)
	})

	t.Run("should return the ErrStartDateMustBeGTDueDate", func(t *testing.T) {
		startDate := time.Now().Format(time.DateTime)
		dueDate := time.Now().Add(time.Minute * -10).Format(time.DateTime)
		_, _, err := validateDates(startDate, dueDate)
		assert.ErrorIs(t, err, ErrStartDateMustBeGTDueDate)
	})

	t.Run("should nil if no error happens", func(t *testing.T) {
		_, _, err := validateDates(validStartDate, validDueDate)
		assert.NoError(t, err)
	})
}
