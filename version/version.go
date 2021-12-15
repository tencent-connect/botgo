// Package version sdk 版本声明。
package version

import (
	"fmt"
)

const (
	// version sdk 版本
	version = "v0.0.1"
	sdkName = "BotGoSDK"
)

// String 输出版本号
func String() string {
	return fmt.Sprintf("%s/%s", sdkName, version)
}
