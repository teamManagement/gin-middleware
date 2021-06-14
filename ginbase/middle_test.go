package ginbase

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestHttpResultUtil(t *testing.T) {
	engine := gin.Default()
	wrapperEngine := HttpResultUtil().RegistryCommonMiddle(engine)
	wrapperEngine.GET("ping", func(ctx *gin.Context) interface{} {
		return "pong"
	})
	wrapperEngine.GET("err", func(ctx *gin.Context) interface{} {
		panic("错误测试")
	})
	wrapperEngine.GET("errCode", func(ctx *gin.Context) interface{} {
		ctx.Set("resCode", 500)
		panic("错误了, 错误码 500")
	})
	_ = engine.Run(":8082")
}
