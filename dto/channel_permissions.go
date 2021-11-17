package dto

// ChannelPermissions 子频道权限
type ChannelPermissions struct {
	ChannelID   string `json:"channel_id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	Permissions string `json:"permissions,omitempty"`
}

// UpdateChannelPermissions 修改子频道权限参数
type UpdateChannelPermissions struct {
	Add    string `json:"add,omitempty"`
	Remove string `json:"remove,omitempty"`
}
