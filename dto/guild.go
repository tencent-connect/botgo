package dto

// Guild 频道结构定义
type Guild struct {
	// 频道ID（与客户端上看到的频道ID不同）
	ID string `json:"id"`
	// 频道名称
	Name string `json:"name"`
	// 频道头像
	Icon string `json:"icon"`
	// 拥有者ID
	OwnerID string `json:"owner_id"`
	// 是否为拥有者
	IsOwner bool `json:"owner"`
	// 成员数量
	MemberCount int `json:"member_count"`
	// 最大成员数目
	MaxMembers int64 `json:"max_members"`
	// 频道描述
	Desc string `json:"description"`
	// 当前用户加入群的时间
	// 此字段只在GUILD_CREATE事件中使用
	JoinedAt Timestamp `json:"joined_at"`
	// 频道列表
	Channels []*Channel `json:"channels"`
	// 游戏绑定公会区服ID
	UnionWorldID string `json:"union_world_id"`
	// 游戏绑定公会/战队ID
	UnionOrgID string `json:"union_org_id"`
	// 操作人
	OpUserID string `json:"op_user_id,omitempty"`
}
