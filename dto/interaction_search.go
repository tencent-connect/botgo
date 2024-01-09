package dto

// SearchInputResolved 搜索类型的输入数据
type SearchInputResolved struct {
	Keyword string `json:"keyword,omitempty"`
}

// SearchRsp 搜索返回数据
type SearchRsp struct {
	Layouts []SearchLayout `json:"layouts"`
}

// SearchLayout 搜索结果的布局
type SearchLayout struct {
	LayoutType LayoutType
	ActionType ActionType
	Title      string
	Records    []SearchRecord
}

// LayoutType 布局类型
type LayoutType uint32

const (
	// LayoutTypeImageText 左图右文
	LayoutTypeImageText LayoutType = 0
)

// ActionType 每行数据的点击行为
type ActionType uint32

const (
	// ActionTypeSendARK 发送 ark 消息
	ActionTypeSendARK ActionType = 0
)

// SearchRecord 每一条搜索结果
type SearchRecord struct {
	Cover string `json:"cover"`
	Title string `json:"title"`
	Tips  string `json:"tips"`
	URL   string `json:"url"`
}

// Resolved 通用的互动反馈
type Resolved struct {
	Keyword     string `json:"keyword"`
	UserID      string `json:"user_id"`
	Request     string `json:"request"`
	MessageID   string `json:"message_id"`
	MemberNick  string `json:"member_nick"`
	ButtonData  string `json:"button_data"`
	ButtonID    string `json:"button_id"`
	FeatureID   string `json:"feature_id"`
	FeedbackOpt string `json:"feedback_opt"` // 智能体反馈选项，LIKE点赞，UNLIKE点踩
	Checked     int32  `json:"checked"`      // 智能体反馈选项是否选中
}
