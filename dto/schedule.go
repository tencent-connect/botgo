package dto

// Schedule 日程对象
type Schedule struct {
	ID             string  `json:"id,omitempty"`
	Name           string  `json:"name,omitempty"`
	Description    string  `json:"description,omitempty"`
	StartTimestamp string  `json:"start_timestamp,omitempty"`
	EndTimestamp   string  `json:"end_timestamp,omitempty"`
	JumpChannelID  string  `json:"jump_channel_id,omitempty"`
	RemindType     string  `json:"remind_type,omitempty"`
	Creator        *Member `json:"creator,omitempty"`
}

// ScheduleWrapper 创建、修改日程的中间对象
type ScheduleWrapper struct {
	Schedule *Schedule `json:"schedule,omitempty"`
}
