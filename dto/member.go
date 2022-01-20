package dto

// Member 群成员
type Member struct {
	GuildID  string    `json:"guild_id"`
	JoinedAt Timestamp `json:"joined_at"`
	Nick     string    `json:"nick"`
	User     *User     `json:"user"`
	Roles    []string  `json:"roles"`
}
