package model

import "github.com/gin-gonic/gin"

func NewSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"message": "success",
		"data":    data,
	})
}
