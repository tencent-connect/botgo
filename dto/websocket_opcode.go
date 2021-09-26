package dto

// OPCode websocket op 码
type OPCode int

// WS OPCode
const (
	WSDispatchEvent OPCode = iota
	WSHeartbeat
	WSIdentity
	_ // Presence Update
	_ // Voice State Update
	_
	WSResume
	WSReconnect
	_ // Request Guild Members
	WSInvalidSession
	WSHello
	WSHeartbeatAck
)
