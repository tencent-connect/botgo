// Package errs 是 SDK 里面的错误类型的集合，同时封装了 SDK 专用的错误类型。
package errs

import (
	"fmt"
)

var (
	// ErrNeedReConnect reconnect
	ErrNeedReConnect = New(CodeNeedReConnect, "need reconnect")
	// ErrInvalidSession 无效的 session
	ErrInvalidSession = New(CodeInvalidSession, "invalid session")
	// ErrURLInvalid ws ap url 异常
	ErrURLInvalid = New(CodeURLInvalid, "ws ap url is invalid")
	// ErrNotFoundOpenAPI 未找到对应版本的openapi实现
	ErrNotFoundOpenAPI = New(CodeNotFoundOpenAPI, "not found openapi version")
	// ErrSessionLimit session 数量受到限制
	ErrSessionLimit = New(CodeSessionLimit, "session num limit")
)

// sdk 错误码
const (
	CodeNeedReConnect = 9000 + iota
	CodeInvalidSession
	CodeURLInvalid
	CodeNotFoundOpenAPI
	CodeSessionLimit
	CodeConnCloseErr // 关闭连接错误码，收拢 websocket close error
)

// Err sdk err
type Err struct {
	code  int
	text  string
	trace string // 错误追踪ID，可用于向平台反馈问题
}

// New 创建一个新错误
func New(code int, text string, trace ...string) error {
	err := &Err{
		code: code,
		text: text,
	}
	if len(trace) > 0 {
		err.trace = trace[0]
	}
	return err
}

// Error 将错误转换为 sdk 的错误类型
func Error(err error) *Err {
	if e, ok := err.(*Err); ok {
		return e
	}
	return &Err{
		code: 9999,
		text: err.Error(),
	}
}

func (e Err) Error() string {
	return fmt.Sprintf("code:%v, text:%v, traceID:%s", e.code, e.text, e.trace)
}

// Code 获取错误码
func (e Err) Code() int {
	return e.code
}

// Text 获取错误信息
func (e Err) Text() string {
	return e.text
}

// Trace 获取错误追踪ID
func (e Err) Trace() string {
	return e.trace
}
