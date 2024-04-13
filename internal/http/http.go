package http

import (
	"github.com/gin-gonic/gin"
)

//go:generate mockgen -destination mocks/controller_mock.go -package mocks . Controller
type Controller interface {
	CreateRoutes(group *gin.RouterGroup)
}

func NewEngine() *gin.Engine {
	engine := gin.Default()
	return engine
}

func StartRoutes(engine *gin.Engine, controllers []Controller) *gin.RouterGroup {
	group := engine.Group("/api")
	for _, controller := range controllers {
		controller.CreateRoutes(group)
	}

	return group
}
