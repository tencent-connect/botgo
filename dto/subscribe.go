package dto

// SubscribeMessageStatusData 订阅消息模板授权数据
type SubscribeMessageStatusData struct {
	GroupOpenid string                       `json:"group_openid"` // 群openid，如果是群订阅消息这里有值
	Openid      string                       `json:"openid"`       // 用户openid，用户订阅消息取这个值
	Result      []SubscribeMsgTemplateResult `json:"result"`       // 授权操作结果
}

// SubscribeMsgTemplateResult 订阅模板授权操作结果
type SubscribeMsgTemplateResult struct {
	TemplateID       int    `json:"template_id"`        // 官方模板id
	CustomTemplateID string `json:"custom_template_id"` // 自定义模板id
	Op               uint32 `json:"op"`                 // 模板授权操作 1-允许 2-拒绝
	SubscribeID      string `json:"subscribe_id"`       // 订阅id
	UpdateTs         uint64 `json:"update_ts"`          // 订阅状态更新时间戳
}
