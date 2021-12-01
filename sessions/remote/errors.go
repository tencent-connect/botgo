package remote

import (
	"errors"
)

var (
	// ErrGotLockFailed 抢锁失败
	ErrGotLockFailed = errors.New("compete for init sessions failed, wait to consume session")
	// ErrSessionMarshalFailed 从redis中读取session后解析失败
	ErrSessionMarshalFailed = errors.New("session marshal failed")
	// ErrProduceFailed 生产session失败
	ErrProduceFailed = errors.New("produce session failed")
	// ErrorNotOk redis 写失败
	ErrorNotOk = errors.New("redis write not ok")
)
