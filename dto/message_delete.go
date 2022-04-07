package dto

// MessageDelete 消息删除结构体定义
type MessageDelete struct {
	// 消息
	Message Message `json:"message"`
	// 操作用户
	OpUser User `json:"op_user"`
}
