package websocket

var (
	// ClientImpl websocket 实现
	ClientImpl WebSocket
)

// Register 注册 websocket 实现
func Register(ws WebSocket) {
	ClientImpl = ws
}
