// Package message 内提供了用于支撑处理消息对象的工具和方法。
package message

import (
	"fmt"
	"regexp"
	"strings"
)

// 用于过滤 at 结构的正则
var atRE = regexp.MustCompile(`<@!\d+>`)

// 用于过滤用户发送消息中的空格符号，\u00A0 是 &nbsp; 的 unicode 编码，某些 mac/pc 版本，连续多个空格的时候会转换成这个符号发送到后台
const spaceCharSet = " \u00A0"

// CMD 一个简单的指令结构
type CMD struct {
	Cmd     string
	Content string
}

// ETLInput 清理输出
//  - 去掉@结构
//  - trim
func ETLInput(input string) string {
	etlData := string(atRE.ReplaceAll([]byte(input), []byte("")))
	etlData = strings.Trim(etlData, spaceCharSet)
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

// Emoji emoji 内嵌格式，参考 https://bot.q.qq.com/wiki/develop/api/openapi/emoji/model.html
// 只支持 type = 1 的系统表情
func Emoji(id int) string {
	return fmt.Sprintf("<emoji:%d>", id)
}

// ParseCommand 解析命令，支持 `{cmd} {content}` 的命令格式
func ParseCommand(input string) *CMD {
	input = ETLInput(input)
	s := strings.Split(input, " ")
	if len(s) < 2 {
		return &CMD{
			Cmd:     strings.Trim(input, spaceCharSet),
			Content: "",
		}
	}
	return &CMD{
		Cmd:     strings.Trim(s[0], spaceCharSet),
		Content: strings.Join(s[1:], " "),
	}
}
