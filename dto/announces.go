package dto

// Announces 公告对象
type Announces struct {
	// 频道 ID
	GuildID string `json:"guild_id"`
	// 子频道 ID
	ChannelID string `json:"channel_id"`
	// 用来创建公告的消息 ID
	MessageID string `json:"message_id"`
	// 公告类别 0:成员公告，1:欢迎公告，默认为成员公告
	AnnouncesType uint32 `json:"announces_type"`
	// 推荐子频道详情数组
	RecommendChannels []RecommendChannel `json:"recommend_channels,omitempty"`
}

// ChannelAnnouncesToCreate 创建子频道公告结构体定义
type ChannelAnnouncesToCreate struct {
	MessageID string `json:"message_id"` // 用来创建公告的消息ID
}

// GuildAnnouncesToCreate 创建频道全局公告结构体定义
type GuildAnnouncesToCreate struct {
	ChannelID         string             `json:"channel_id"`         // 用来创建公告的子频道 ID
	MessageID         string             `json:"message_id"`         // 用来创建公告的消息 ID
	AnnouncesType     uint32             `json:"announces_type"`     // 公告类别 0:成员公告，1:欢迎公告，默认为成员公告
	RecommendChannels []RecommendChannel `json:"recommend_channels"` // 推荐子频道详情列表
}

// RecommendChannel 推荐子频道详情
type RecommendChannel struct {
	ChannelID string `json:"channel_id"` // 子频道 ID
	Introduce string `json:"introduce"`  // 推荐语
}
