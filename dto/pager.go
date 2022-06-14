package dto

// Pager 分页器接口需要实现将对象转换为分页参数的方法
type Pager interface {
	// QueryParams 转换为 query 参数
	QueryParams() map[string]string
}

// GuildMembersPager 分页器
type GuildMembersPager struct {
	After string `json:"after"` // 读此id之后的数据，如果是第一次请求填0，默认为0
	Limit string `json:"limit"` // 分页大小，1-1000，默认是1
}

// QueryParams 转换为 query 参数
func (g *GuildMembersPager) QueryParams() map[string]string {
	query := make(map[string]string)
	if g.Limit != "" {
		query["limit"] = g.Limit
	}
	if g.After != "" {
		query["after"] = g.After
	}
	return query
}

// GuildRoleMembersPager 分页器
type GuildRoleMembersPager struct {
	StartIndex string `json:"start_index"`
	Limit      string `json:"limit"`
}

// QueryParams 转换为 query 参数
func (g *GuildRoleMembersPager) QueryParams() map[string]string {
	query := make(map[string]string)
	if g.Limit != "" {
		query["limit"] = g.Limit
	}
	if g.StartIndex != "" {
		query["start_index"] = g.StartIndex
	}
	return query
}

// GuildPager 分页器
type GuildPager struct {
	Before string `json:"before"` // 读此id之前的数据
	After  string `json:"after"`  // 读此id之后的数据
	Limit  string `json:"limit"`  // 分页大小，1-100，默认是 100
}

// QueryParams 转换为 query 参数
func (g *GuildPager) QueryParams() map[string]string {
	query := make(map[string]string)
	if g.Limit != "" {
		query["limit"] = g.Limit
	}
	if g.After != "" {
		query["after"] = g.After
	}
	// 优先 after
	if g.After == "" && g.Before != "" {
		query["before"] = g.Before
	}
	return query
}

// MessagesPager 消息分页
type MessagesPager struct {
	Type  MessagePagerType // 拉取类型
	ID    string           // 消息ID
	Limit string           `json:"limit"` // 最大 20
}

// QueryParams 转换为 query 参数
func (m *MessagesPager) QueryParams() map[string]string {
	query := make(map[string]string)
	if m.Limit != "" {
		query["limit"] = m.Limit
	}
	if m.Type != "" && m.ID != "" {
		query[string(m.Type)] = m.ID
	}
	return query
}

// MessageReactionPager 分页器
type MessageReactionPager struct {
	Cookie string `json:"cookie"` // 分页游标
	Limit  string `json:"limit"`  // 分页大小，1-1000，默认是20
}

// QueryParams 转换为 query 参数
func (g *MessageReactionPager) QueryParams() map[string]string {
	query := make(map[string]string)
	if g.Limit != "" {
		query["limit"] = g.Limit
	}
	if g.Cookie != "" {
		query["cookie"] = g.Cookie
	}
	return query
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
