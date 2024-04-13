package http

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"todo-app/internal/http/mocks"
)

func TestNewEngine(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should create new engine with default config", func(t *testing.T) {
		engine := NewEngine()
		assert.NotNil(t, engine)
	})
}

func TestStart(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should create the controllers routes", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		controller := mocks.NewMockController(ctrl)
		controller.EXPECT().CreateRoutes(gomock.Any())

		engine := gin.Default()
		group := StartRoutes(engine, []Controller{controller})
		assert.NotNil(t, group)
	})
}
