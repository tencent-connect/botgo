package dto

// PinsMessage 精华消息对象
type PinsMessage struct {
	// 频道 ID
	GuildID string `json:"guild_id"`
	// 子频道 ID
	ChannelID string `json:"channel_id"`
	// 消息 ID 数组
	MessageIDs []string `json:"message_ids"`
}
