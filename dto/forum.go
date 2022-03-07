package dto

// Thread 主题事件主体内容
type Thread struct {
	GuildID    string     `json:"guild_id"`
	ChannelID  string     `json:"channel_id"`
	AuthorID   string     `json:"author_id"`
	ThreadInfo ThreadInfo `json:"thread_info"`
}

// ThreadInfo 主题信息
type ThreadInfo struct {
	ThreadID string `json:"thread_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	DateTime string `json:"date_time"`
}

// Post 帖子事件主体内容
type Post struct {
	GuildID   string   `json:"guild_id"`
	ChannelID string   `json:"channel_id"`
	AuthorID  string   `json:"author_id"`
	PostInfo  PostInfo `json:"post_info"`
}

// PostInfo 帖子内容
type PostInfo struct {
	ThreadID string `json:"thread_id"`
	PostID   string `json:"post_id"`
	Content  string `json:"content"`
	DateTime string `json:"date_time"`
}

// Reply 回复事件主体内容
type Reply struct {
	GuildID   string    `json:"guild_id"`
	ChannelID string    `json:"channel_id"`
	AuthorID  string    `json:"author_id"`
	ReplyInfo ReplyInfo `json:"reply_info"`
}

// ReplyInfo 回复内容
type ReplyInfo struct {
	ThreadID string `json:"thread_id"`
	PostID   string `json:"post_id"`
	ReplyID  string `json:"reply_id"`
	Content  string `json:"content"`
	DateTime string `json:"date_time"`
}

// ForumAuditResult 帖子审核事件主体内容
type ForumAuditResult struct {
	TaskID      string `json:"task_id"`
	GuildID     string `json:"guild_id"`
	ChannelID   string `json:"channel_id"`
	AuthorID    string `json:"author_id"`
	ThreadID    string `json:"thread_id"`
	PostID      string `json:"post_id"`
	ReplyID     string `json:"reply_id"`
	PublishType uint32 `json:"type"`
	Result      uint32 `json:"result"`
	ErrMsg      string `json:"err_msg"`
	DateTime    string `json:"date_time"`
}
