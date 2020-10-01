package service

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	limitgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// GetRouter constructs a Gin Router, binds all routes, and then returns the resulting router for later use.
func GetRouter() *gin.Engine {
	router := gin.Default()
	router.ForwardedByClientIP = true
	if gin.Mode() == gin.DebugMode {
		router.Use(cors.New(cors.Config{
			AllowOrigins: []string{"http://localhost:8000"},
			AllowMethods: []string{"POST"},
			AllowHeaders: []string{"content-type"},
		}))
	}

	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  5,
	}
	store := memory.NewStore()
	limit := limitgin.NewMiddleware(limiter.New(store, rate))

	router.Use(limit)

	router.POST("/process", _ProcessPost)

	return router
}
