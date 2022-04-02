// Package websocket SDK 需要实现的 websocket 定义。
package websocket

import (
	"runtime"
	"syscall"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/log"
)

var (
	// ClientImpl websocket 实现
	ClientImpl WebSocket
	// ResumeSignal 用于强制 resume 连接的信号量
	ResumeSignal syscall.Signal
)

// Register 注册 websocket 实现
func Register(ws WebSocket) {
	ClientImpl = ws
}

// RegisterResumeSignal 注册用于通知 client 将连接进行 resume 的信号
func RegisterResumeSignal(signal syscall.Signal) {
	ResumeSignal = signal
}

// PanicBufLen Panic 堆栈大小
var PanicBufLen = 1024

// PanicHandler 处理websocket场景的 panic ，打印堆栈
func PanicHandler(e interface{}, session *dto.Session) {
	buf := make([]byte, PanicBufLen)
	buf = buf[:runtime.Stack(buf, false)]
	log.Errorf("[PANIC]%s\n%v\n%s\n", session, e, buf)
}

// RegisterHandlers 兼容老版本的注册方式
func RegisterHandlers(handlers ...interface{}) dto.Intent {
	return event.RegisterHandlers(handlers...)
}
