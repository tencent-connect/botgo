package dto

// GuildRoles 频道用户组列表返回
type GuildRoles struct {
	GuildID  string  `json:"guild_id"`
	Roles    []*Role `json:"roles"`
	NumLimit string  `json:"role_num_limit"`
}

// Role 频道身份组
type Role struct {
	ID          RoleID `json:"id,omitempty"`
	Name        string `json:"name"`
	Color       uint32 `json:"color"`
	Hoist       uint32 `json:"hoist"`
	MemberCount uint32 `json:"number,omitempty"`       // 不会被修改，创建接口修改
	MemberLimit uint32 `json:"member_limit,omitempty"` // 不会被修改，创建接口修改
}

// DefaultColor 用户组默认颜色值
const DefaultColor = 4278245297

// RoleID 用户组ID
type RoleID string

// UpdateRoleInfo 身份组可更改数据
type UpdateRoleInfo struct {
	Name  string `json:"name"`
	Color uint32 `json:"color"`
	Hoist uint32 `json:"hoist"`
}

// UpdateRoleFilter 身份组可更改数据，修改的
type UpdateRoleFilter struct {
	Name  uint32 `json:"name"`
	Color uint32 `json:"color"`
	Hoist uint32 `json:"hoist"`
}

// UpdateRole role 更新请求承载
type UpdateRole struct {
	GuildID string            `json:"guild_id"`
	Filter  *UpdateRoleFilter `json:"filter"`
	Update  *Role             `json:"info"`
}

// UpdateResult 创建，删除等行为的返回
type UpdateResult struct {
	RoleID  `json:"role_id"`
	GuildID string `json:"guild_id"`
	Role    *Role  `json:"role"`
}

// MemberAddRoleBody 增加子频道管理员身份组时附加内容
type MemberAddRoleBody struct {
	// 子频道对象
	Channel *Channel `json:"channel"`
}
