package service

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// GetRouter constructs a Gin Router, binds all routes, and then returns the resulting router for later use.
func GetRouter() *gin.Engine {
	router := gin.Default()
	if gin.Mode() == gin.DebugMode {
		router.Use(cors.New(cors.Config{
			AllowOrigins: []string{"http://localhost:8000"},
			AllowMethods: []string{"POST"},
			AllowHeaders: []string{"content-type"},
		}))
	}

	router.POST("/process", _ProcessPost)

	return router
}
