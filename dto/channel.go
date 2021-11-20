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
	ChannelTypeLive        = 10000 + iota // 直播子频道
	ChannelTypeApplication                // 应用子频道
)

// ChannelSubType 子频道子类型定义
type ChannelSubType int

// 子频道子类型定义
const (
	ChannelSubTypeChat     ChannelSubType = iota // 闲聊，默认子类型
	ChannelSubTypeNotice                         // 公告
	ChannelSubTypeGuide                          // 攻略
	ChannelSubTypeTeamGame                       // 开黑
)

// Channel 频道结构定义
type Channel struct {
	// 频道ID
	ID string `json:"id"`
	// 群ID
	GuildID string `json:"guild_id"`
	ChannelValueObject
}

// ChannelValueObject 频道的值对象部分
type ChannelValueObject struct {
	// 频道名称
	Name string `json:"name"`
	// 频道类型
	Type ChannelType `json:"type"`
	// 排序位置
	Position int64 `json:"position"`
	// 父频道的ID
	ParentID string `json:"parent_id"`
	// 拥有者ID
	OwnerID string `json:"owner_id"`
	// 子频道子类型
	SubType ChannelSubType `json:"sub_type"`
}
