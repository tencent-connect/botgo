package dto

// UpdateGuildMute 更新频道相关禁言的Body参数
type UpdateGuildMute struct {
	// 禁言截止时间戳，单位秒
	MuteEndTimestamp string `json:"mute_end_timestamp,omitempty"`
	// 禁言多少秒（两个字段二选一，默认以mute_end_timstamp为准）
	MuteSeconds string `json:"mute_seconds,omitempty"`
}
