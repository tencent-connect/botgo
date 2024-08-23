package dto

// EnterAIO 进入aio的事件
type EnterAIO struct {
	UserOpenid string `json:"user_openid,omitempty"` // 用户openid
	FromSource string `json:"from_source,omitempty"` // 进入aio的来源
}
