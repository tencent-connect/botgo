// Package errs 是 SDK 里面的错误类型的集合，同时封装了 SDK 专用的错误类型。
package errs

import (
	"fmt"
)

var (
	// ErrNeedReConnect reconnect
	ErrNeedReConnect = New(CodeNeedReConnect, "need reconnect")
	// ErrInvalidSession 无效的 session
	ErrInvalidSession = New(CodeConnCloseCantResume, "invalid session")
	// ErrURLInvalid ws ap url 异常
	ErrURLInvalid = New(CodeConnCloseCantIdentify, "ws ap url is invalid")
	// ErrSessionLimit session 数量受到限制
	ErrSessionLimit = New(CodeConnCloseCantIdentify, "session num limit")

	// ErrNotFoundOpenAPI 未找到对应版本的openapi实现
	ErrNotFoundOpenAPI = New(CodeNotFoundOpenAPI, "not found openapi version")
	// ErrPagerIsNil 分页器为空
	ErrPagerIsNil = New(CodePagerIsNil, "pager is nil")
)

// sdk 错误码
const (
	CodeNeedReConnect = 9000
	// CodeInvalidSession 无效的的 session id 请重新连接
	CodeInvalidSession  = 9001
	CodeURLInvalid      = 9002
	CodeNotFoundOpenAPI = 9003
	CodeSessionLimit    = 9004
	// CodeConnCloseCantResume 关闭连接错误码，收拢 websocket close error，不允许 resume
	CodeConnCloseCantResume = 9005
	// CodeConnCloseCantIdentify 不允许连接的关闭连接错误，比如机器人被封禁
	CodeConnCloseCantIdentify = 9006
	// CodePagerIsNil 分页器为空
	CodePagerIsNil = 9007
)

// websocket错误码
const (
	WSCodeBackendUnknownError         = 4000 // 未知错误
	WSCodeBackendUnknownOpCode        = 4001 // 请求中的opcode非法
	WSCodeBackendDecodeError          = 4002 // 请求发送的数据解析失败，数据格式有问题
	WSCodeBackendNotAuthenticate      = 4003 // 尚未进行身份校验
	WSCodeBackendAuthenticationFail   = 4004 // 身份校验失败
	WSCodeBackendAlreadyAuthenticate  = 4005 // 重复的身份校验
	WSCodeBackendSessionNoLongerValid = 4006 // session已失效
	WSCodeBackendInvalidSeq           = 4007 // 非法的序号
	WSCodeBackendRateLimit            = 4008 // 请求超频
	WSCodeBackendSessionTimeOut       = 4009 // session过期，需要重连
	WSCodeBackendInvalidShard         = 4010 // 非法的shard
	WSCodeBackendShardingRequired     = 4011 // 当前shard承载数据过多，需要重连已重新分配shard
	WSCodeBackendInvalidAPIVersion    = 4012 // 非法的api版本号
	WSCodeBackendInvalidIntents       = 4013 // 非法的intents参数
	WSCodeBackendDisallowdIntents     = 4014 // 使用了未授权使用的intents
	WSCodeBackendBotOffline           = 4914 // 机器人已经被下架
	WSCodeBackendBotBanned            = 4915 // 机器人被封禁
)

// openapi错误码
const (
	APICodeTokenExpireOrNotExist = 11244 // token过期或者不存在
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

// Error 输出错误信息
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
