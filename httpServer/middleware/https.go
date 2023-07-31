package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RedirectToHTTPS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.TLS == nil {
			host := ctx.Request.Host
			ctx.Redirect(http.StatusMovedPermanently, "https://"+host+ctx.Request.RequestURI)
			return
		}
		ctx.Next()
	}
}
