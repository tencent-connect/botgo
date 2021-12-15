// Package openapi 声明了 sdk 所使用的 openapi 接口。
package openapi

import (
	"net/http"
	"sync"
)

// VersionMapping openapi 版本管理
var VersionMapping = map[APIVersion]OpenAPI{}

// DefaultImpl 默认 openapi 实现
var DefaultImpl OpenAPI

var (
	versionMapLock = sync.RWMutex{}
	once           sync.Once
)

// 这些状态码不会当做错误处理
// 未排除 201,202 : 用于提示创建异步任务成功，所以不屏蔽错误
var successStatusSet = map[int]bool{
	http.StatusOK:        true,
	http.StatusNoContent: true,
}

// Register 注册 openapi 的实现，如果默认实现为空，则将第一个注册的设置为默认实现
func Register(version APIVersion, api OpenAPI) {
	versionMapLock.Lock()
	VersionMapping[version] = api
	setDefaultOnce(api)
	versionMapLock.Unlock()
}

// IsSuccessStatus 是否是成功的状态码
func IsSuccessStatus(code int) bool {
	if _, ok := successStatusSet[code]; ok {
		return true
	}
	return false
}

func setDefaultOnce(api OpenAPI) {
	once.Do(func() {
		if DefaultImpl == nil {
			DefaultImpl = api
		}
	})
}
