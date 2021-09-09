package dto

import (
	"fmt"
	"strings"
	"time"
)

// Duration 支持能够直接配置中解析出来 time.Duration 类型的数据
// 需要实现对应类型的 Unmarshaler 接口
type Duration time.Duration

// UnmarshalJSON 实现json的解析接口
func (d *Duration) UnmarshalJSON(bytes []byte) error {
	var s = strings.Trim(string(bytes), "\"'")
	t, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("failed to parse '%s' to time.Duration: %v", s, err)
	}

	*d = Duration(t)
	return nil
}
