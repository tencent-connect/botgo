package dto

// UpdateGuildMute 更新频道相关禁言的Body参数
type UpdateGuildMute struct {
	// 禁言截止时间戳，单位秒
	MuteEndTimestamp string `json:"mute_end_timestamp,omitempty"`
	// 禁言多少秒（两个字段二选一，默认以mute_end_timstamp为准）
	MuteSeconds string `json:"mute_seconds,omitempty"`
	// 批量禁言的成员列表（全员禁言时不填写该字段）
	UserIDs []string `json:"user_ids,omitempty"`
}

// UpdateGuildMuteResponse 批量禁言的回参
type UpdateGuildMuteResponse struct {
	// 批量禁言成功的成员列表
	UserIDs []string `json:"user_ids,omitempty"`
}
