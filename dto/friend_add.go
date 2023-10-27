package dto

type C2CFriendData struct {
	OpenId    string `json:"openid"`
	Timestamp int    `json:"timestamp"` // 添加/删除机器人好友时间戳
	Nick      string `json:"nick"`      // 待事件链路补充
	Avatar    string `json:"avatar"`    // 待事件链路补充
}
