package ginmiddleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UseVerifyUserAgent(targetUserAgent string) gin.HandlerFunc {
	return UseVerifyUserAgentWithVerifyFn(func(userAgent string) bool {
		return userAgent == targetUserAgent
	})
}

func UseVerifyUserAgentWithVerifyFn(fn func(userAgent string) bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userAgent := ctx.GetHeader("User-Agent")
		if !fn(userAgent) {
			ctx.Abort()
			ctx.Status(http.StatusMethodNotAllowed)
			return
		}
		ctx.Next()
	}
}
