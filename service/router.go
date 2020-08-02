package service

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// GetRouter constructs a Gin Router, binds all routes, and then returns the resulting router for later use.
func GetRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true, // TODO: Lock this down.
		AllowMethods:    []string{"POST"},
	}))

	router.POST("/process", _ProcessPost)

	return router
}
