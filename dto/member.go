package dto

// Member 群成员
type Member struct {
	GuildID  string    `json:"guild_id"`
	JoinedAt Timestamp `json:"joined_at"`
	Nick     string    `json:"nick"`
	User     *User     `json:"user"`
	Roles    []string  `json:"roles"`
	OpUserID string    `json:"op_user_id,omitempty"`
}

// DeleteHistoryMsgDay 消息撤回天数
type DeleteHistoryMsgDay = int

// 支持的消息撤回天数，除这些天数之外，传递其他值将不会撤回任何消息
const (
	NoDelete                              = 0  // 不删除任何消息
	DeleteThreeDays   DeleteHistoryMsgDay = 3  // 3天
	DeleteSevenDays   DeleteHistoryMsgDay = 7  // 7天
	DeleteFifteenDays DeleteHistoryMsgDay = 15 // 15天
	DeleteThirtyDays  DeleteHistoryMsgDay = 30 // 30天
	DeleteAll         DeleteHistoryMsgDay = -1 // 删除所有消息
)

// MemberDeleteOpts 删除成员额外参数
type MemberDeleteOpts struct {
	AddBlackList         bool                `json:"add_blacklist"`
	DeleteHistoryMsgDays DeleteHistoryMsgDay `json:"delete_history_msg_days"`
}

// MemberDeleteOption 删除成员选项
type MemberDeleteOption func(*MemberDeleteOpts)

// WithAddBlackList 将当前成员同时添加到频道黑名单中
func WithAddBlackList(b bool) MemberDeleteOption {
	return func(o *MemberDeleteOpts) {
		o.AddBlackList = b
	}
}

// WithDeleteHistoryMsg 删除成员时同时撤回消息
func WithDeleteHistoryMsg(days DeleteHistoryMsgDay) MemberDeleteOption {
	return func(o *MemberDeleteOpts) {
		o.DeleteHistoryMsgDays = days
	}
}
