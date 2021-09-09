package dto

// ChannelType 频道类型定义
type ChannelType int

// 子频道类型定义
const (
	ChannelTypeText ChannelType = iota
	_
	ChannelTypeVoice
	_
	ChannelTypeCategory
	ChannelTypeLive        // 直播子频道
	ChannelTypeApplication // 应用子频道
)

// Channel 频道结构定义
type Channel struct {
	// 频道ID
	ID string `json:"id"`
	// 群ID
	GuildID string `json:"guild_id"`
	// 频道名称
	Name string `json:"name"`
	// 频道类型
	Type ChannelType `json:"type"`
	// 排序位置
	Position int `json:"position"`
	// 父频道的ID
	ParentID string `json:"parent_id"`
	// 拥有者ID
	OwnerID string `json:"owner_id"`
}
