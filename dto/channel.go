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
	ChannelTypeForum                      // 论坛子频道
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

// ChannelPrivateType 频道可见性类型定义
type ChannelPrivateType int

// 频道可见性类型定义
const (
	ChannelPrivateTypePublic         ChannelPrivateType = iota // 公开频道
	ChannelPrivateTypeOnlyAdmin                                // 群主管理员可见
	ChannelPrivateTypeAdminAndMember                           // 群主管理员+指定成员
)

// SpeakPermissionType 发言权限类型定义
type SpeakPermissionType int

// 发言权限类型定义
const (
	SpeakPermissionTypePublic         SpeakPermissionType = iota + 1 // 公开发言权限
	SpeakPermissionTypeAdminAndMember                                // 指定成员可发言
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
	Name string `json:"name,omitempty"`
	// 频道类型
	Type ChannelType `json:"type,omitempty"`
	// 排序位置
	Position int64 `json:"position,omitempty"`
	// 父频道的ID
	ParentID string `json:"parent_id,omitempty"`
	// 拥有者ID
	OwnerID string `json:"owner_id,omitempty"`
	// 子频道子类型
	SubType ChannelSubType `json:"sub_type,omitempty"`
	// 子频道可见性类型
	PrivateType ChannelPrivateType `json:"private_type,omitempty"`
	// 创建私密子频道的时候，同时带上 userID，能够将这些成员添加为私密子频道的成员
	// 注意：只有创建私密子频道的时候才支持这个参数
	PrivateUserIDs []string `json:"private_user_ids,omitempty"`
	// 发言权限
	SpeakPermission SpeakPermissionType `json:"speak_permission,omitempty"`
	// 应用子频道的应用ID，仅应用子频道有效，定义请参考
	// [文档](https://bot.q.qq.com/wiki/develop/api/openapi/channel/model.html)
	ApplicationID string `json:"application_id,omitempty"`
	// 机器人在此频道上拥有的权限, 定义请参考
	// [文档](https://bot.q.qq.com/wiki/develop/api/openapi/channel_permissions/model.html#permissions)
	Permissions string `json:"permissions,omitempty"`
	// 操作人
	OpUserID string `json:"op_user_id,omitempty"`
}
