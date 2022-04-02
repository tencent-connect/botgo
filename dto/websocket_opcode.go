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
	HTTPCallbackAck
)

// opMeans op 对应的含义字符串标识
var opMeans = map[OPCode]string{
	WSDispatchEvent:  "Event",
	WSHeartbeat:      "Heartbeat",
	WSIdentity:       "Identity",
	WSResume:         "Resume",
	WSReconnect:      "Reconnect",
	WSInvalidSession: "InvalidSession",
	WSHello:          "Hello",
	WSHeartbeatAck:   "HeartbeatAck",
}

// OPMeans 返回 op 含义
func OPMeans(op OPCode) string {
	means, ok := opMeans[op]
	if !ok {
		means = "unknown"
	}
	return means
}
