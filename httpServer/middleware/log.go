package middleware

import (
	"cook-robot-middle-platform-go/logger"
	"github.com/gin-gonic/gin"
)

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 根据请求路径或其他条件判断是否要跳过打印日志
		if c.Request.URL.Path == "/api/v1/controller/fetchStatus" {
			// 禁止打印该请求的日志

			logger.Log.Println("123")
		}

		c.Next()
	}
}
