// Package message 内提供了用于支撑处理消息对象的工具和方法。
package message

import (
	"fmt"
	"regexp"
	"strings"
)

// 用于过滤 at 结构的正则
var atRE = regexp.MustCompile(`<@!\d+>`)

// ETLInput 清理输出
//  - 去掉@结构
//  - trim
func ETLInput(input string) string {
	etlData := string(atRE.ReplaceAll([]byte(input), []byte("")))
	etlData = strings.Trim(etlData, " ")
	return etlData
}

// MentionUser 返回 at 用户的内嵌格式
// https://bot.q.qq.com/wiki/develop/api/openapi/message/message_format.html
func MentionUser(userID string) string {
	return fmt.Sprintf("<@%s>", userID)
}

// MentionAllUser 返回 at all 的内嵌格式
func MentionAllUser() string {
	return "@everyone"
}

// MentionChannel 提到子频道的格式
func MentionChannel(channelID string) string {
	return fmt.Sprintf("<#%s>", channelID)
}
