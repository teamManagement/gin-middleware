package ginbase

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var httpResultUtil = &httpResultUtils{}

func HttpResultUtil() *httpResultUtils {
	return httpResultUtil
}

type HttpResult struct {
	Code   string      `json:"code,omitempty"`
	Error  bool        `json:"error,omitempty"`
	Msg    string      `json:"msg,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

type ServiceFun func(ctx *gin.Context) interface{}

type httpResultUtils struct {
}

type WrapperGinEngine struct {
	*gin.Engine
	httpResultUtils *httpResultUtils
}

func (w *WrapperGinEngine) GET(relativePath string, fun ServiceFun) *WrapperGinEngine {
	w.Engine.GET(relativePath, w.httpResultUtils.WrapperResponseHandle(fun))
	return w
}

func (w *WrapperGinEngine) POST(relativePath string, fun ServiceFun) *WrapperGinEngine {
	w.Engine.POST(relativePath, w.httpResultUtils.WrapperResponseHandle(fun))
	return w
}

func (w *WrapperGinEngine) DELETE(relativePath string, fun ServiceFun) *WrapperGinEngine {
	w.Engine.DELETE(relativePath, w.httpResultUtils.WrapperResponseHandle(fun))
	return w
}

func (w *WrapperGinEngine) PATCH(relativePath string, fun ServiceFun) *WrapperGinEngine {
	w.Engine.PATCH(relativePath, w.httpResultUtils.WrapperResponseHandle(fun))
	return w
}

func (w *WrapperGinEngine) PUT(relativePath string, fun ServiceFun) *WrapperGinEngine {
	w.Engine.PUT(relativePath, w.httpResultUtils.WrapperResponseHandle(fun))
	return w
}

func (w *WrapperGinEngine) OPTIONS(relativePath string, fun ServiceFun) *WrapperGinEngine {
	w.Engine.OPTIONS(relativePath, w.httpResultUtils.WrapperResponseHandle(fun))
	return w
}

func (w *WrapperGinEngine) HEAD(relativePath string, fun ServiceFun) *WrapperGinEngine {
	w.Engine.HEAD(relativePath, w.httpResultUtils.WrapperResponseHandle(fun))
	return w
}

// NewErrorHttpResultWithMsg 新创建一个错误的httpResult根据错误信息
func (h *httpResultUtils) NewErrorHttpResultWithMsg(str string) *HttpResult {
	return &HttpResult{
		Code:  "1",
		Msg:   str,
		Error: true,
	}
}

// NewErrorHttpResultWithCodeAndMsg 创建一个错误的httpResult根据错误信息
func (h httpResultUtils) NewErrorHttpResultWithCodeAndMsg(code, msg string) *HttpResult {
	return &HttpResult{
		Code:  code,
		Msg:   msg,
		Error: true,
	}
}

// NewSuccessHttpResult 创建一个成功的httpResult
func (h *httpResultUtils) NewSuccessHttpResult() *HttpResult {
	return &HttpResult{
		Code:  "0",
		Error: false,
	}
}

// NewSuccessHttpResultWithResult 创建正确信息根据结果
func (h *httpResultUtils) NewSuccessHttpResultWithResult(result interface{}) *HttpResult {
	return &HttpResult{
		Code:   "0",
		Error:  false,
		Result: result,
	}
}

// WrapperResponseHandle 包装结果
func (h *httpResultUtils) WrapperResponseHandle(fn ServiceFun) gin.HandlerFunc {
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
			context.JSON(code, h.NewSuccessHttpResultWithResult(v))

		}
	}
}

// WrapperErrorHandle 包装错误回调器
func (h *httpResultUtils) MiddleErrorHandle() gin.HandlerFunc {
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
				logrus.Error(e)
				context.JSON(resCodeNum, h.NewErrorHttpResultWithMsg("未知异常"))
			case string:
				context.JSON(resCodeNum, h.NewErrorHttpResultWithMsg(v))
			case HttpResult:
				context.JSON(resCodeNum, v)
			default:
				context.JSON(resCodeNum, v)
			}
		}()
		context.Next()
	}
}

func (h *httpResultUtils) RegistryCommonMiddle(engine *gin.Engine) *WrapperGinEngine {
	engine.Use(h.MiddleErrorHandle())
	return &WrapperGinEngine{
		engine,
		HttpResultUtil(),
	}
}
