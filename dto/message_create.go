package dto

import "github.com/tencent-connect/botgo/dto/keyboard"

// SendType 消息类型
type SendType int

const (
	Text      SendType = 1 // Text 文字消息
	RichMedia SendType = 2 // RichMedia 富媒体类消息
)

// APIMessage 消息结构接口
type APIMessage interface {
	GetEventID() string
	GetSendType() SendType
}

// RichMediaMessage 富媒体消息
// 注意：直接使用srv_send_msg=tre时会占用主动消息频率，且多媒体文件不能复用，建议先上传，然后再使用消息发送类型7进行发送
type RichMediaMessage struct {
	EventID    string `json:"event_id,omitempty"`     // 已经废弃：要回复的事件id, 逻辑同MsgID
	FileType   uint64 `json:"file_type,omitempty"`    // 业务类型，图片，文件，语音，视频 文件类型，取值:1图片,2视频,3语音(目前语音只支持silk格式)
	URL        string `json:"url,omitempty"`          // 需发送的富媒体文件，HTTP或者HTTPS链接
	SrvSendMsg bool   `json:"srv_send_msg,omitempty"` // 为true时会直接发送到群/C2C，且会占用主动消息频率, 为false为上传富媒体文件
	Content    string `json:"content,omitempty"`
	MsgSeq     int64  `json:"msg_seq,omitempty"` // 机器人对于回复一个msg_id或者event_id的消息序号，指定后根据这个字段和msg_id或者event_id进行去重
}

// GetEventID 事件ID
func (msg RichMediaMessage) GetEventID() string {
	return ""
}

// GetSendType 消息类型
func (msg RichMediaMessage) GetSendType() SendType {
	return RichMedia
}

// MessageType 消息类型
type MessageType int

const (
	TextMsg        MessageType = 0 // 文字消息
	MarkdownMsg    MessageType = 2 // md 消息
	ArkMsg         MessageType = 3 // ark消息类型
	EmbedMsg       MessageType = 4 // EMBED消息
	ATMsg          MessageType = 5 // @消息
	InputNotifyMsg MessageType = 6 // 输入状态消息
	RichMediaMsg   MessageType = 7 // 富媒体消息（图片，视频等）
)

// MessageToCreate 发送消息结构体定义
type MessageToCreate struct {
	Content string      `json:"content,omitempty"`
	MsgType MessageType `json:"msg_type,omitempty"` //消息类型: 0:文字消息, 2: md消息
	Embed   *Embed      `json:"embed,omitempty"`
	Ark     *Ark        `json:"ark,omitempty"`
	Image   string      `json:"image,omitempty"`
	// 要回复的消息id，为空是主动消息，公域机器人会异步审核，不为空是被动消息，公域机器人会校验语料
	MsgID            string                    `json:"msg_id,omitempty"`
	MessageReference *MessageReference         `json:"message_reference,omitempty"`
	Markdown         *Markdown                 `json:"markdown,omitempty"`
	Keyboard         *keyboard.MessageKeyboard `json:"keyboard,omitempty"`        // 消息按钮组件
	EventID          string                    `json:"event_id,omitempty"`        // 要回复的事件id, 逻辑同MsgID
	Timestamp        int64                     `json:"timestamp,omitempty"`       //TODO delete this
	MsgSeq           uint32                    `json:"msg_seq,omitempty"`         // 机器人对于回复一个msg_id或者event_id的消息序号，指定后根据这个字段和msg_id或者event_id进行去重
	SubscribeID      string                    `json:"subscribe_id,omitempty"`    // 订阅id，发送订阅消息时使用
	InputNotify      *InputNotify              `json:"input_notify,omitempty"`    // 输入状态状态信息
	Media            *MediaInfo                `json:"media,omitempty"`           // 富媒体信息
	PromptKeyboard   *PromptKeyboard           `json:"prompt_keyboard,omitempty"` // 消息扩展信息
	ActionButton     *ActionButton             `json:"action_button,omitempty"`   // 消息操作结构
	Stream           *Stream                   `json:"stream,omitempty"`          // 流式消息信息
	FeatureID        uint                      `json:"feature_id,omitempty"`      // 控制消息发送
}

// Stream 流式消息信息
type Stream struct {
	State int32  `json:"state,omitempty"` // 流式消息状态 1正文生成中，10：正文生成结束， 11：引志消息生成中， 20：引导消息生成结束。
	ID    string `json:"id,omitempty"`    // 流式消息ID，流式消息第一条不用填写，第二条需要填写第一个分片返回的msgID.
	Index int32  `json:"index,omitempty"` // 流式消息的序号， 从1开始
	Reset bool   `json:"reset,omitempty"` // 重新生成流式消息标记，此参数只能使用于流式消息分片还没有发送完成时，reset时Index需要从0开始，需要填写流式ID。
}

// PromptKeyboard 交互区操作
type PromptKeyboard struct {
	Keyboard *keyboard.MessageKeyboard `json:"keyboard,omitempty"` // 消息按钮组件
}

// ActionButton 消息操作按钮
type ActionButton struct {
	TemplateID   int32  `json:"template_id,omitempty"`   // 消息操作栏模块ID，与下面具体具体按钮二选一填写。待废弃字段！！！
	CallbackData string `json:"callback_data,omitempty"` // 用户操作时会回调通过回调事件给到button_data中， 最长不超过128个字符。
	Feedback     bool   `json:"feedback,omitempty"`      // 反馈按钮（赞踩按钮）
	TTS          bool   `json:"tts,omitempty"`           // TTS语音播放按钮
	ReGenerate   bool   `json:"re_generate,omitempty"`   // 重新生成按钮
	StopGenerate bool   `json:"stop_generate,omitempty"` // 停止生成按钮
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
	TemplateID       int               `json:"template_id"`        // 模版 id
	CustomTemplateID string            `json:"custom_template_id"` // 自定义模板id
	Params           []*MarkdownParams `json:"params"`             // 模版参数
	Content          string            `json:"content"`            // 原生 markdown
	Style            *MarkdownStyle    `json:"style"`              // markdown样式
	ProcessMsg       string            `json:"process_msg"`        // markdown引导消息
}

// MarkdownStyle markdown 样式
type MarkdownStyle struct {
	MainFontSize string `json:"main_font_size"` // 正文字体大小 small middle large
	Layout       string `json:"layout"`         // hide_avatar_and_center 隐藏头像并居中
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

// InputNotify 输入状态结构
type InputNotify struct {
	InputType   int   `json:"input_type,omitempty"`   // 类型 1: "对方正在输入...", 2: 取消展示"]
	InputSecond int32 `json:"input_second,omitempty"` // 当input_type大于0时有效, 代码状态持续多长时间.
}

// MediaInfo 富媒体信息
type MediaInfo struct {
	FileInfo []byte `json:"file_info,omitempty"` // 富媒体文件信息，通过上传接口取得
}
