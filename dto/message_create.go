package dto

import "github.com/tencent-connect/botgo/dto/keyboard"

// SendType 消息类型
type SendType int

const (
	Text      SendType = 1 // 文字消息
	RichMedia SendType = 2 // 富媒体类消息
)

// APIMessage 消息结构接口
type APIMessage interface {
	GetEventID() string
	GetSendType() SendType
}

// RichMediaMessage 富媒体消息
type RichMediaMessage struct {
	EventID    string `json:"event_id,omitempty"`  // 要回复的事件id, 逻辑同MsgID
	FileType   uint64 `json:"file_type,omitempty"` // 业务类型，图片，文件，语音，视频 文件类型，取值:1图片,2视频,3语音(目前语音只支持silk格式)
	URL        string `json:"url,omitempty"`
	SrvSendMsg bool   `json:"srv_send_msg,omitempty"`
	Content    string `json:"content,omitempty"`
}

// GetEventID 事件ID
func (msg RichMediaMessage) GetEventID() string {
	return msg.EventID
}

// GetSendType 消息类型
func (msg RichMediaMessage) GetSendType() SendType {
	return RichMedia
}

// MessageToCreate 发送消息结构体定义
type MessageToCreate struct {
	Content string `json:"content,omitempty"`
	MsgType int    `json:"msg_type,omitempty"` //消息类型: 0:文字消息, 2: md消息
	Embed   *Embed `json:"embed,omitempty"`
	Ark     *Ark   `json:"ark,omitempty"`
	Image   string `json:"image,omitempty"`
	// 要回复的消息id，为空是主动消息，公域机器人会异步审核，不为空是被动消息，公域机器人会校验语料
	MsgID            string                    `json:"msg_id,omitempty"`
	MessageReference *MessageReference         `json:"message_reference,omitempty"`
	Markdown         *Markdown                 `json:"markdown,omitempty"`
	Keyboard         *keyboard.MessageKeyboard `json:"keyboard,omitempty"`  // 消息按钮组件
	EventID          string                    `json:"event_id,omitempty"`  // 要回复的事件id, 逻辑同MsgID
	Timestamp        int64                     `json:"timestamp,omitempty"` //TODO delete this
}

// GetEventID 事件ID
func (msg MessageToCreate) GetEventID() string {
	return msg.EventID
}

// GetSendType 消息类型
func (msg MessageToCreate) GetSendType() SendType {
	return Text
}

// MessageReference 引用消息
type MessageReference struct {
	MessageID             string `json:"message_id"`               // 消息 id
	IgnoreGetMessageError bool   `json:"ignore_get_message_error"` // 是否忽律获取消息失败错误
}

// GetEventID 事件ID
func (msg MessageReference) GetEventID() string {
	return msg.MessageID
}

// GetSendType 消息类型
func (msg MessageReference) GetSendType() SendType {
	return Text
}

// Markdown markdown 消息
type Markdown struct {
	TemplateID int               `json:"template_id"` // 模版 id
	Params     []*MarkdownParams `json:"params"`      // 模版参数
	Content    string            `json:"content"`     // 原生 markdown
}

// MarkdownParams markdown 模版参数 键值对
type MarkdownParams struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// SettingGuideToCreate 发送引导消息的结构体
type SettingGuideToCreate struct {
	Content      string        `json:"content,omitempty"`       // 频道内发引导消息可以带@
	SettingGuide *SettingGuide `json:"setting_guide,omitempty"` // 设置引导
}

// SettingGuide 设置引导
type SettingGuide struct {
	// 频道ID, 当通过私信发送设置引导消息时，需要指定guild_id
	GuildID string `json:"guild_id"`
}
