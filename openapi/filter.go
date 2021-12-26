package openapi

import (
	"net/http"
	"sync"
)

// 提供一组过滤器支持，开发者可以通过请求过滤器和返回过滤器，实现模调上报，耗时监控等能力。

// HTTPFilter 请求过滤器
type HTTPFilter func(req *http.Request, response *http.Response) error

var (
	filterLock         = sync.RWMutex{}
	reqFilterChainSet  = map[string]HTTPFilter{}
	reqFilterChains    []string
	respFilterChainSet = map[string]HTTPFilter{}
	respFilterChains   []string
)

// RegisterReqFilter 注册请求过滤器
func RegisterReqFilter(name string, filter HTTPFilter) {
	if _, ok := reqFilterChainSet[name]; ok {
		return
	}
	filterLock.Lock()
	defer filterLock.Unlock()
	reqFilterChainSet[name] = filter
	reqFilterChains = append(reqFilterChains, name)
}

// RegisterRespFilter 注册返回过滤器
func RegisterRespFilter(name string, filter HTTPFilter) {
	if _, ok := respFilterChainSet[name]; ok {
		return
	}
	filterLock.Lock()
	defer filterLock.Unlock()
	respFilterChainSet[name] = filter
	respFilterChains = append(respFilterChains, name)
}

// DoReqFilterChains 按照注册顺序执行请求过滤器
func DoReqFilterChains(req *http.Request, resp *http.Response) error {
	for _, name := range reqFilterChains {
		if _, ok := reqFilterChainSet[name]; !ok {
			continue
		}
		if err := reqFilterChainSet[name](req, resp); err != nil {
			return err
		}
	}
	return nil
}

// DoRespFilterChains 按照注册顺序执行返回过滤器
func DoRespFilterChains(req *http.Request, resp *http.Response) error {
	for _, name := range respFilterChains {
		if _, ok := respFilterChainSet[name]; !ok {
			continue
		}
		if err := respFilterChainSet[name](req, resp); err != nil {
			return err
		}
	}
	return nil
}
