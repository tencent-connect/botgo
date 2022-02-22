package dto

// MessageToCreate 发送消息结构体定义
type MessageToCreate struct {
	Content          string            `json:"content,omitempty"`
	Embed            *Embed            `json:"embed,omitempty"`
	Ark              *Ark              `json:"ark,omitempty"`
	Image            string            `json:"image,omitempty"`
	MsgID            string            `json:"msg_id,omitempty"` // 要回复的消息id，为空是主动消息，公域机器人会异步审核，不为空是被动消息，公域机器人会校验语料
	MessageReference *MessageReference `json:"message_reference,omitempty"`
	Markdown         *Markdown         `json:"markdown,omitempty"`
	Keyboard         string            `json:"keyboard,omitempty"` // 内嵌键盘
}

// MessageReference 引用消息
type MessageReference struct {
	MessageID             string `json:"message_id"`               // 消息 id
	IgnoreGetMessageError bool   `json:"ignore_get_message_error"` // 是否忽律获取消息失败错误
}

// Markdown markdown 消息
type Markdown struct {
	TemplateID int               `json:"template_id"`
	Params     []*MarkdownParams `json:"params"`
}

// MarkdownParams markdown 模版参数 键值对
type MarkdownParams struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// Keyboard 内嵌键盘
type Keyboard struct {
	ID string `json:"id"`
}
