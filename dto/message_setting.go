package dto

// MessageSetting 消息频率设置信息
type MessageSetting struct {
	DisableCreateDm   bool     `json:"disable_create_dm,omitempty"`
	DisablePushMsg    bool     `json:"disable_push_msg,omitempty"`
	ChannelIDs        []string `json:"channel_ids,omitempty"`
	ChannelPushMaxNum int      `json:"channel_push_max_num,omitempty"`
}
