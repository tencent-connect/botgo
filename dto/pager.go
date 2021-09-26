package dto

// GuildMembersPager 分页器
type GuildMembersPager struct {
	After string `json:"after"` // 上一次回包中最大的ID， 如果是第一次请求填0，默认为0
	Limit string `json:"limit"` // 分页大小，1-1000，默认是1
}

// MessagesPager 消息分页
type MessagesPager struct {
	Type  MessagePagerType // 拉取类型
	ID    string           // 消息ID
	Limit string           `json:"limit"` // 最大 20
}

// MessagePagerType 消息翻页拉取方式
type MessagePagerType string

const (
	// MPTAround 拉取消息ID上下的消息
	MPTAround MessagePagerType = "around"
	// MPTBefore 拉取消息ID之前的消息
	MPTBefore MessagePagerType = "before"
	// MPTAfter 拉取消息ID之后的消息
	MPTAfter MessagePagerType = "after"
)
