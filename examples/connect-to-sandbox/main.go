package main

import (
	"context"
	"log"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/token"
)

func main() {
	ctx := context.Background()
	// 加载 appid 和 token
	botToken := token.New(token.TypeBot)
	if err := botToken.LoadFromConfig("config.yaml"); err != nil {
		log.Fatalln(err)
	}
	// 初始化 openapi，使用 NewSandboxOpenAPI 请求到沙箱环境
	api := botgo.NewSandboxOpenAPI(botToken).WithTimeout(3 * time.Second)
	// 获取 websocket 信息，如果 api 是请求到沙箱环境的，则获取到沙箱环境的 ws 地址
	// websocket 的链接，以及事件处理，请参考其他 examples
	wsInfo, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(wsInfo.URL)
}
