package dto

// DirectMessage 私信结构定义，一个 DirectMessage 为两个用户之间的一个私信频道，简写为 DM
type DirectMessage struct {
	// 频道ID
	GuildID string `json:"guild_id"`
	// 子频道id
	ChannelID string `json:"channel_id"`
	// 私信频道创建的时间戳
	CreateTime string `json:"create_time"`
}

// DirectMessageToCreate 创建私信频道的结构体定义
type DirectMessageToCreate struct {
	// 频道ID
	SourceGuildID string `json:"source_guild_id"`
	// 用户ID
	RecipientID string `json:"recipient_id"`
}
