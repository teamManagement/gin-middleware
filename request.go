package ginmiddleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func UseNotFoundHandle(engine *gin.Engine) {
	engine.NoRoute(notFoundHandle)
	engine.NoMethod(notFoundHandle)
}

func notFoundHandle(ctx *gin.Context) {
	ctx.JSON(404, NewErrorHttpResultWithCodeAndMsg("404", fmt.Sprintf("Not Found %s %s", ctx.Request.Method, ctx.Request.URL.String())))
}
