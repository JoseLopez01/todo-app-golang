package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"todo-app/todo"
	"todo-app/todo/dtos"
)

type TodosController struct {
	service todo.Service
}

func NewTodosController(service todo.Service) *TodosController {
	return &TodosController{
		service: service,
	}
}

func (t *TodosController) CreateRoutes(base *gin.RouterGroup) {
	group := base.Group("/todos/:email")
	group.POST("", t.Create)
	group.GET("", t.GetAll)
	group.GET("/:id", t.GetByID)
	group.DELETE(":id", t.Delete)
	group.PUT(":id", t.Update)
}

func (t *TodosController) Create(ctx *gin.Context) {
	var dto dtos.CreateTodo
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	email := ctx.Param("email")
	response, err := t.service.Create(ctx, email, dto)
	if err != nil {
		code := getStatusCode(err)
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": response})
}

func (t *TodosController) GetAll(ctx *gin.Context) {
	email := ctx.Param("email")
	response, err := t.service.GetAll(ctx, email)
	if err != nil {
		code := getStatusCode(err)
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func (t *TodosController) GetByID(ctx *gin.Context) {
	email := ctx.Param("email")
	id := ctx.Param("id")
	response, err := t.service.GetByID(ctx, email, id)
	if err != nil {
		code := getStatusCode(err)
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func (t *TodosController) Delete(ctx *gin.Context) {
	email := ctx.Param("email")
	id := ctx.Param("id")

	if err := t.service.Delete(ctx, email, id); err != nil {
		code := getStatusCode(err)
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (t *TodosController) Update(ctx *gin.Context) {
	var dto dtos.UpdateTodo
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	email := ctx.Param("email")
	id := ctx.Param("id")
	response, err := t.service.Update(ctx, email, id, dto)
	if err != nil {
		code := getStatusCode(err)
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func getStatusCode(err error) int {
	if errors.Is(err, todo.ErrInvalidID) || errors.Is(err, todo.ErrInvalidDueDate) || errors.Is(err, todo.ErrInvalidStartDate) || errors.Is(err, todo.ErrStartDateMustBeGTDueDate) {
		return http.StatusBadRequest
	}

	if errors.Is(err, todo.ErrTodoIsCompleted) {
		return http.StatusConflict
	}

	return http.StatusInternalServerError
}
