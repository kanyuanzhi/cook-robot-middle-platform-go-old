package model

import "github.com/gin-gonic/gin"

func NewSuccessResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    data,
	})
}

func NewFailResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, gin.H{
		"message": "error",
		"data":    data,
	})
}
