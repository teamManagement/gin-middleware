package ginmiddleware

import "fmt"

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
func NewErrorHttpResultWithMsg(format string, args ...any) *HttpResult {
	return &HttpResult{
		Code:  "1",
		Msg:   fmt.Sprintf(format, args...),
		Error: true,
	}
}

// NewErrorHttpResultWithCodeAndMsg 创建一个错误的httpResult根据错误信息
func NewErrorHttpResultWithCodeAndMsg(code, format string, args ...any) *HttpResult {
	return &HttpResult{
		Code:  code,
		Msg:   fmt.Sprintf(format, args...),
		Error: true,
	}
}

// NewSuccessHttpResult 创建一个成功的httpResult
func NewSuccessHttpResult() *HttpResult {
	return &HttpResult{
		Code:  "0",
		Error: false,
	}
}
