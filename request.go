package ginmiddleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// UseNotFoundHandle 使用没有找到路由的拦截器
func UseNotFoundHandle(engine *gin.Engine) {
	engine.NoRoute(notFoundHandle)
	engine.NoMethod(notFoundHandle)
}

// notFoundHandle 没有找到路由的拦截器
func notFoundHandle(ctx *gin.Context) {
	ctx.JSON(404, NewErrorHttpResultWithCodeAndMsg("404", fmt.Sprintf("Not Found %s %s", ctx.Request.Method, ctx.Request.URL.String())))
}
