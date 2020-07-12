package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type _ProcessBody struct {
	Code string `json:"code" binding:"required"`
}

func _ProcessPost(c *gin.Context) {
	var body _ProcessBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(body)
}
