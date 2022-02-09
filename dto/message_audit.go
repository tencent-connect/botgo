package dto

// MessageAudit 消息审核结构体定义
type MessageAudit struct {
	// 审核 ID
	AuditID string `json:"audit_id"`
	// 消息 ID
	MessageID string `json:"message_id"`
	// 频道 ID
	GuildID string `json:"guild_id"`
	// 子频道 ID
	ChannelID string `json:"channel_id"`
	// 审核时间
	AuditTime string `json:"audit_time"`
	// 创建时间
	CreateTime string `json:"create_time"`
	// 子频道 seq，用于消息间的排序，seq 在同一子频道中按从先到后的顺序递增，不同的子频道之前消息无法排序
	SeqInChannel string `json:"seq_in_channel"`
}
