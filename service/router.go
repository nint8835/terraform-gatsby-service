package service

import "github.com/gin-gonic/gin"

// GetRouter constructs a Gin Router, binds all routes, and then returns the resulting router for later use.
func GetRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/process", _ProcessPost)

	return router
}
