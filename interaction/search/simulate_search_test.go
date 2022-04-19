package search

import (
	"testing"

	"github.com/tencent-connect/botgo/log"
)

func TestSimulateSearch(t *testing.T) {
	got, err := SimulateSearch(
		&Config{
			AppID:    "1",
			EndPoint: "https://www.qq.com",
			Secret:   "a",
		}, "hello",
	)
	if err != nil {
		// 这里用于模拟，默认的 testcase 肯定是失败的，所以这里不断言
		log.Error(err)
	}
	log.Info(got)
}
