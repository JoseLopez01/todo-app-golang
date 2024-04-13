package main

import (
	"go.uber.org/fx"

	"todo-app/config"
	"todo-app/internal/http"
	"todo-app/internal/storage"
	"todo-app/todo"
)

func main() {
	app := fx.New(
		fx.Supply(config.AppConfig),
		storage.Module,
		http.Module,
		todo.Module,
	)

	app.Run()
}
