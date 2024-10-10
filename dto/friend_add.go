package dto

// C2CFriendData c2c 好友事件信息
type C2CFriendData struct {
	OpenID    string `json:"openid"`
	Timestamp int    `json:"timestamp"` // 添加/删除机器人好友时间戳
	Nick      string `json:"nick"`      // 待事件链路补充
	Avatar    string `json:"avatar"`    // 待事件链路补充
}
