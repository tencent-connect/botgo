package dto

// Member 群成员
type Member struct {
	GuildID  string    `json:"guild_id"`
	JoinedAt Timestamp `json:"joined_at"`
	Nick     string    `json:"nick"`
	User     *User     `json:"user"`
	Roles    []string  `json:"roles"`
}

// MemberDeleteOpts 删除成员额外参数
type MemberDeleteOpts struct {
	AddBlackList bool `json:"add_blacklist"`
}

// MemberDeleteOption 删除成员选项
type MemberDeleteOption func(*MemberDeleteOpts)

// WithAddBlackList 将当前成员同时添加到频道黑名单中
func WithAddBlackList(b bool) MemberDeleteOption {
	return func(o *MemberDeleteOpts) {
		o.AddBlackList = b
	}
}
