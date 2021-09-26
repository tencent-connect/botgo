package openapi

import (
	"fmt"
)

type (
	// APIVersion 接口版本
	APIVersion = uint32
)

// 接口版本，后续增加版本直接直接增加新的常量
const (
	APIv1 APIVersion = 1 + iota
	_
)

// APIVersionString 返回version的字符串格式 ex: v1
func APIVersionString(version APIVersion) string {
	return fmt.Sprintf("v%v", version)
}
