package dto

// AudioStatus 音频状态
type AudioStatus uint32

// 音频状态
const (
	AudioStatusStart = iota
	AudioStatusPause
	AudioStatusResume
	AudioStatusStop
)

// AudioControl 音频控制对象
type AudioControl struct {
	URL    string      `json:"audio_url"`
	Text   string      `json:"text"`
	Status AudioStatus `json:"status"`
}

// AudioAction 音频动作
type AudioAction struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	URL       string `json:"audio_url"`
	Text      string `json:"text"`
}
