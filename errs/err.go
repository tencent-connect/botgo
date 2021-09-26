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
	code int
	text string
}

// New 创建一个新错误
func New(code int, text string) error {
	return &Err{
		code: code,
		text: text,
	}
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
	return fmt.Sprintf("code:%v, text:%v", e.code, e.text)
}

// Code 获取错误码
func (e Err) Code() int {
	return e.code
}

// Text 获取错误信息
func (e Err) Text() string {
	return e.text
}
