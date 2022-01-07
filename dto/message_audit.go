package dto

// MessageAudit 消息审核结构体定义
type MessageAudit struct {
	AuditID    string `json:"audit_id"`
	MessageID  string `json:"message_id"`
	GuildID    string `json:"guild_id"`
	ChannelID  string `json:"channel_id"`
	AuditTime  string `json:"audit_time"`
	CreateTime string `json:"create_time"`
}
