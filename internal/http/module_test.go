package http

import (
	"github.com/gin-gonic/gin"
	"testing"
	"todo-app/config"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/mock/gomock"

	"todo-app/todo"
	"todo-app/todo/mocks"
)

func TestAsController(t *testing.T) {
	t.Run("should return a not nil controller", func(t *testing.T) {
		response := AsController(func() {})
		assert.NotNil(t, response)
	})
}

func TestModule(t *testing.T) {
	t.Run("should return create the module", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		app := fxtest.New(
			t,
			fx.Supply(config.Config{}),
			fx.Provide(
				fx.Annotate(
					func() todo.Service {
						return mocks.NewMockService(ctrl)
					},
					fx.As(new(todo.Service)),
				),
				fx.Annotate(
					func(engine *gin.Engine) bool {
						return engine != nil
					},
					fx.ParamTags(engineTag),
				),
			),
			fx.Invoke(func(engine bool) {
				assert.True(t, engine)
			}),
			Module,
		)
		defer app.RequireStart().RequireStop()
	})
}
