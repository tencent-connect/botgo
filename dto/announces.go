package dto

// Announces 公告对象
type Announces struct {
	// 频道ID
	GuildID string `json:"guild_id"`
	// 子频道id
	ChannelID string `json:"channel_id"`
	// 用来创建公告的消息ID
	MessageID string `json:"message_id"`
}

// ChannelAnnouncesToCreate 创建子频道公告结构体定义
type ChannelAnnouncesToCreate struct {
	MessageID string `json:"message_id"` // 用来创建公告的消息ID
}

// GuildAnnouncesToCreate 创建频道全局公告结构体定义
type GuildAnnouncesToCreate struct {
	ChannelID string `json:"channel_id"` // 用来创建公告的子频道ID
	MessageID string `json:"message_id"` // 用来创建公告的消息ID
}
