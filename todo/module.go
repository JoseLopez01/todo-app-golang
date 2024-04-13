package todo

import "go.uber.org/fx"

var Module = fx.Module(
	"todo-module",
	fx.Provide(
		fx.Private,
		fx.Annotate(
			NewRedisRepository,
			fx.As(new(Repository)),
		),
	),
	fx.Provide(
		fx.Annotate(
			NewTodosService,
			fx.As(new(Service)),
		),
	),
)
