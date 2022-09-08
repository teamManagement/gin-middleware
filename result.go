package ginmiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-base-lib/logs"
)

type WrapperGinEngine struct {
	*gin.Engine
}

type ServiceFun func(ctx *gin.Context) interface{}

func (w *WrapperGinEngine) GET(relativePath string, fun ServiceFun) *WrapperGinEngine {
	w.Engine.GET(relativePath, WrapperResponseHandle(fun))
	return w
}

func (w *WrapperGinEngine) POST(relativePath string, fun ServiceFun) *WrapperGinEngine {
	w.Engine.POST(relativePath, WrapperResponseHandle(fun))
	return w
}

func (w *WrapperGinEngine) DELETE(relativePath string, fun ServiceFun) *WrapperGinEngine {
	w.Engine.DELETE(relativePath, WrapperResponseHandle(fun))
	return w
}

func (w *WrapperGinEngine) PATCH(relativePath string, fun ServiceFun) *WrapperGinEngine {
	w.Engine.PATCH(relativePath, WrapperResponseHandle(fun))
	return w
}

func (w *WrapperGinEngine) PUT(relativePath string, fun ServiceFun) *WrapperGinEngine {
	w.Engine.PUT(relativePath, WrapperResponseHandle(fun))
	return w
}

func (w *WrapperGinEngine) OPTIONS(relativePath string, fun ServiceFun) *WrapperGinEngine {
	w.Engine.OPTIONS(relativePath, WrapperResponseHandle(fun))
	return w
}

func (w *WrapperGinEngine) HEAD(relativePath string, fun ServiceFun) *WrapperGinEngine {
	w.Engine.HEAD(relativePath, WrapperResponseHandle(fun))
	return w
}

// WrapperResponseHandle 包装结果
func WrapperResponseHandle(fn ServiceFun) gin.HandlerFunc {
	return func(context *gin.Context) {
		result := fn(context)
		resCode, exists := context.Get("resCode")
		if !exists {
			resCode = 0
		}

		code, ok := resCode.(int)
		if !ok {
			code = 200
		}

		switch v := result.(type) {
		case HttpResult:
			context.JSON(code, v)
		case error:
			panic(v)
		default:
			context.JSON(code, NewSuccessHttpResultWithResult(v))

		}
	}
}

// HttpResult http响应结果封装
type HttpResult struct {
	// Code 消息代码
	Code string `json:"code,omitempty"`
	// Error 消息是否有错
	Error bool `json:"error,omitempty"`
	// Msg 错误消息
	Msg string `json:"msg,omitempty"`
	// Result 返回的对象
	Result interface{} `json:"result,omitempty"`
}

// NewSuccessHttpResultWithResult 创建正确信息根据结果
func NewSuccessHttpResultWithResult(result interface{}) *HttpResult {
	return &HttpResult{
		Code:   "0",
		Error:  false,
		Result: result,
	}
}

// NewErrorHttpResultWithMsg 新创建一个错误的httpResult根据错误信息
func NewErrorHttpResultWithMsg(str string) *HttpResult {
	return &HttpResult{
		Code:  "1",
		Msg:   str,
		Error: true,
	}
}

func UseRecover2HttpResult() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			e := recover()
			if e == nil {
				return
			}
			resCode, exists := context.Get("resCode")
			if !exists {
				resCode = 200
			}

			resCodeNum, ok := resCode.(int)
			if !ok {
				resCodeNum = 200
			}

			switch v := e.(type) {
			case error:
				logs.Error(e)
				context.JSON(resCodeNum, NewErrorHttpResultWithMsg("未知异常"))
			case string:
				context.JSON(resCodeNum, NewErrorHttpResultWithMsg(v))
			case HttpResult:
				context.JSON(resCodeNum, v)
			default:
				context.JSON(resCodeNum, v)
			}
		}()
		context.Next()
	}
}
