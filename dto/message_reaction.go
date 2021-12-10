package dto

// ReactionTargetType 表情表态对象类型
type ReactionTargetType = int32

const (
	ReactionTargetTypeMsg     = iota // 消息
	ReactionTargetTypeFeed           // 帖子
	ReactionTargetTypeComment        // 评论
	ReactionTargetTypeReply          // 回复
)

// MessageReaction 表情表态动作
type MessageReaction struct {
	UserId    string         `json:"user_id"`
	ChannelId string         `json:"channel_id"`
	GuildId   string         `json:"guild_id"`
	Target    ReactionTarget `json:"target"`
	Emoji     Emoji          `json:"emoji"`
}

// ReactionTarget 表态对象类型
type ReactionTarget struct {
	Id   string             `json:"id"`
	Type ReactionTargetType `json:"type"`
}
