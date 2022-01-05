// Package websocket SDK 需要实现的 websocket 定义。
package websocket

import (
	"runtime"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/log"
)

var (
	// ClientImpl websocket 实现
	ClientImpl WebSocket
)

// Register 注册 websocket 实现
func Register(ws WebSocket) {
	ClientImpl = ws
}

// PanicBufLen Panic 堆栈大小
var PanicBufLen = 1024

// PanicHandler 处理websocket场景的 panic ，打印堆栈
func PanicHandler(e interface{}, session *dto.Session) {
	buf := make([]byte, PanicBufLen)
	buf = buf[:runtime.Stack(buf, false)]
	log.Errorf("[PANIC]%s\n%v\n%s\n", session, e, buf)
}
