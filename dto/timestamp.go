package dto

import "time"

// Timestamp 时间戳
type Timestamp string

// Time 时间字符串格式转换
func (t Timestamp) Time() (time.Time, error) {
	return time.Parse(time.RFC3339, string(t))
}
