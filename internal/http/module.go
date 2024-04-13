package http

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"todo-app/config"
	"todo-app/internal/http/controllers"
)

const (
	controllersTag = `group:"controllers"`
	engineTag      = `name:"engine"`
)

var Module = fx.Module(
	"http-module",
	fx.Provide(
		StartServer,
		fx.Annotate(NewEngine, fx.ResultTags(engineTag)),
		fx.Annotate(StartRoutes, fx.ParamTags(engineTag, controllersTag)),
		AsController(controllers.NewTodosController),
	),
	fx.Invoke(func(*gin.RouterGroup) {}),
)

func StartServer(configs config.Config, engine *gin.Engine, lc fx.Lifecycle) bool {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := engine.Run(configs.Port); err != nil {
					panic(err)
				}
			}()

			return nil
		},
	})

	return true
}

func AsController(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Controller)),
		fx.ResultTags(controllersTag),
	)
}
