package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"todo-app/todo"
	"todo-app/todo/dtos"
	"todo-app/todo/mocks"
	"todo-app/todo/models"
)

var (
	ctxMatcher   = gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem())
	emailMatcher = gomock.Eq("test@example.com")
	idMatcher    = gomock.Eq("279f4a4e-48dc-4569-83df-8b30ce488599")
)

func TestNewTodosController(t *testing.T) {
	t.Run("should create new controller", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockService(ctrl)
		controller := NewTodosController(service)
		assert.NotNil(t, controller)
	})
}

func TestTodosController_CreateRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("should create new routes", func(t *testing.T) {
		engine := gin.Default()
		group := engine.Group("/api")
		controller := NewTodosController(nil)
		controller.CreateRoutes(group)

		routes := engine.Routes()
		assert.Len(t, routes, 5)
	})
}

func TestTodosController_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	dtoMatcher := gomock.AssignableToTypeOf(dtos.CreateTodo{})

	dto, err := json.Marshal(dtos.CreateTodo{
		Description: "description",
		Name:        "name",
		StartDate:   time.Now().Format(time.DateTime),
		DueDate:     time.Now().Add(time.Minute).Format(time.DateTime),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should return 400 if the request is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockService(ctrl)

		r := gin.Default()

		controller := NewTodosController(service)
		controller.CreateRoutes(r.Group("/api"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/todos/test@test.test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return the service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockService(ctrl)
		service.EXPECT().Create(ctxMatcher, emailMatcher, dtoMatcher).Return(models.Todo{}, fmt.Errorf("error"))

		r := gin.Default()

		controller := NewTodosController(service)
		controller.CreateRoutes(r.Group("/api"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/todos/test@example.com", bytes.NewReader(dto))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return 201 if the todo is created", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockService(ctrl)
		service.EXPECT().Create(ctxMatcher, emailMatcher, dtoMatcher).Return(models.Todo{}, nil)

		r := gin.Default()

		controller := NewTodosController(service)
		controller.CreateRoutes(r.Group("/api"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/todos/test@example.com", bytes.NewReader(dto))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestTodosController_GetAll(t *testing.T) {
	t.Run("should return the service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockService(ctrl)
		service.EXPECT().GetAll(ctxMatcher, emailMatcher).Return(nil, fmt.Errorf("error"))

		r := gin.Default()
		controller := NewTodosController(service)
		controller.CreateRoutes(r.Group("/api"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/todos/test@example.com", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockService(ctrl)
		service.EXPECT().GetAll(ctxMatcher, emailMatcher).Return([]models.Todo{{ID: "279f4a4e-48dc-4569-83df-8b30ce488599"}}, nil)

		r := gin.Default()
		controller := NewTodosController(service)
		controller.CreateRoutes(r.Group("/api"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/todos/test@example.com", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTodosController_GetByID(t *testing.T) {
	t.Run("should return the service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockService(ctrl)
		service.EXPECT().GetByID(ctxMatcher, emailMatcher, idMatcher).Return(models.Todo{}, fmt.Errorf("error"))

		r := gin.Default()
		controller := NewTodosController(service)
		controller.CreateRoutes(r.Group("/api"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/todos/test@example.com/279f4a4e-48dc-4569-83df-8b30ce488599", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockService(ctrl)
		service.EXPECT().GetByID(ctxMatcher, emailMatcher, idMatcher).Return(models.Todo{ID: "id"}, nil)

		r := gin.Default()
		controller := NewTodosController(service)
		controller.CreateRoutes(r.Group("/api"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/todos/test@example.com/279f4a4e-48dc-4569-83df-8b30ce488599", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTodosController_Update(t *testing.T) {
	dtoMatcher := gomock.AssignableToTypeOf(dtos.UpdateTodo{})
	dto, err := json.Marshal(dtos.UpdateTodo{})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should return 400 if the request is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockService(ctrl)

		r := gin.Default()
		controller := NewTodosController(service)
		controller.CreateRoutes(r.Group("/api"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/api/todos/test@example.com/279f4a4e-48dc-4569-83df-8b30ce488599", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return the service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockService(ctrl)
		service.EXPECT().Update(ctxMatcher, emailMatcher, idMatcher, dtoMatcher).Return(models.Todo{}, fmt.Errorf("error"))

		r := gin.Default()
		controller := NewTodosController(service)
		controller.CreateRoutes(r.Group("/api"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/api/todos/test@example.com/279f4a4e-48dc-4569-83df-8b30ce488599", bytes.NewReader(dto))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockService(ctrl)
		service.EXPECT().Update(ctxMatcher, emailMatcher, idMatcher, dtoMatcher).Return(models.Todo{}, nil)

		r := gin.Default()
		controller := NewTodosController(service)
		controller.CreateRoutes(r.Group("/api"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/api/todos/test@example.com/279f4a4e-48dc-4569-83df-8b30ce488599", bytes.NewReader(dto))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTodosController_Delete(t *testing.T) {
	t.Run("should return the service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockService(ctrl)
		service.EXPECT().Delete(ctxMatcher, emailMatcher, idMatcher).Return(fmt.Errorf("error"))

		r := gin.Default()
		controller := NewTodosController(service)
		controller.CreateRoutes(r.Group("/api"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/api/todos/test@example.com/279f4a4e-48dc-4569-83df-8b30ce488599", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return 204", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockService(ctrl)
		service.EXPECT().Delete(ctxMatcher, emailMatcher, idMatcher).Return(nil)

		r := gin.Default()
		controller := NewTodosController(service)
		controller.CreateRoutes(r.Group("/api"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/api/todos/test@example.com/279f4a4e-48dc-4569-83df-8b30ce488599", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func Test_GetStatusCode(t *testing.T) {
	t.Run("should return 400 for user errors", func(t *testing.T) {
		userErrors := []error{
			todo.ErrInvalidID,
			todo.ErrInvalidDueDate,
			todo.ErrInvalidStartDate,
			todo.ErrStartDateMustBeGTDueDate,
		}

		for _, err := range userErrors {
			code := getStatusCode(err)
			assert.Equal(t, http.StatusBadRequest, code)
		}
	})

	t.Run("should return 409", func(t *testing.T) {
		code := getStatusCode(todo.ErrTodoIsCompleted)
		assert.Equal(t, http.StatusConflict, code)
	})

	t.Run("should return 500 for any other error", func(t *testing.T) {
		code := getStatusCode(fmt.Errorf("error"))
		assert.Equal(t, http.StatusInternalServerError, code)
	})
}
