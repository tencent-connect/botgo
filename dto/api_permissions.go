package dto

// APIPermissions API 权限列表对象
type APIPermissions struct {
	APIList []*APIPermission `json:"apis,omitempty"` // API 权限列表
}

// APIPermission API 权限对象
type APIPermission struct {
	Path       string `json:"path,omitempty"`        // API 接口名，例如 /guilds/{guild_id}/members/{user_id}
	Method     string `json:"method,omitempty"`      // 请求方法，例如 GET
	Desc       string `json:"desc,omitempty"`        // API 接口名称，例如 获取频道信
	AuthStatus int    `json:"auth_status,omitempty"` // 授权状态，auth_stats 为 1 时已授权
}

// APIPermissionDemandIdentify  API 权限需求标识对象
type APIPermissionDemandIdentify struct {
	Path   string `json:"path,omitempty"`   // API 接口名，例如 /guilds/{guild_id}/members/{user_id}
	Method string `json:"method,omitempty"` // 请求方法，例如 GET
}

// APIPermissionDemand 接口权限需求对象
type APIPermissionDemand struct {
	GuildID     string                       `json:"guild_id,omitempty"`     // 频道 ID
	ChannelID   string                       `json:"channel_id,omitempty"`   // 子频道 ID
	APIIdentify *APIPermissionDemandIdentify `json:"api_identify,omitempty"` // 权限接口唯一标识
	Title       string                       `json:"title,omitempty"`        // 接口权限链接中的接口权限描述信息
	Desc        string                       `json:"desc,omitempty"`         // 接口权限链接中的机器人可使用功能的描述信息
}

// APIPermissionDemandToCreate 创建频道 API 接口权限授权链接结构体定义
type APIPermissionDemandToCreate struct {
	ChannelID   string                       `json:"channel_id"`             // 子频道 ID
	APIIdentify *APIPermissionDemandIdentify `json:"api_identify,omitempty"` // 接口权限链接中的接口权限描述信息
	Desc        string                       `json:"desc"`                   // 接口权限链接中的机器人可使用功能的描述信息
}
