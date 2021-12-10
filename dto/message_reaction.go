package dto

// ReactionTargetType 表情表态对象类型
type ReactionTargetType = int32

const (
	// ReactionTargetTypeMsg 消息
	ReactionTargetTypeMsg = iota
	// ReactionTargetTypeFeed 帖子
	ReactionTargetTypeFeed
	// ReactionTargetTypeComment 评论
	ReactionTargetTypeComment
	// ReactionTargetTypeReply 回复
	ReactionTargetTypeReply
)

// MessageReaction 表情表态动作
type MessageReaction struct {
	UserID    string         `json:"user_id"`
	ChannelID string         `json:"channel_id"`
	GuildID   string         `json:"guild_id"`
	Target    ReactionTarget `json:"target"`
	Emoji     Emoji          `json:"emoji"`
}

// ReactionTarget 表态对象类型
type ReactionTarget struct {
	ID   string             `json:"id"`
	Type ReactionTargetType `json:"type"`
}
