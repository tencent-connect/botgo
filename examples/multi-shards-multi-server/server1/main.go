package main

import (
	"context"
	"log"
	"time"

	"examples/multi-shards-multi-server/handler"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/sessions/remote"
	"github.com/tencent-connect/botgo/token"
)

func main() {
	ctx := context.Background()

	// 把默认的 session manager 替换为 remote session manager，提供 redis 连接
	// 把集群 key 设置为 "two-shard-two-server"，所有使用相同集群 key 的进程，会一起进行竞争
	botgo.SetSessionManager(
		remote.New(
			handler.GetRedisConn("localhost:6379", 3*time.Second),
			remote.WithClusterKey("two-shard-two-server"),
		),
	)

	// 加载 appid 和 token
	botToken := token.New(token.TypeBot)
	if err := botToken.LoadFromConfig("config.yaml"); err != nil {
		log.Fatalln(err)
	}
	// 初始化 openapi
	api := botgo.NewOpenAPI(botToken).WithTimeout(3 * time.Second)

	// 获取 websocket 信息
	wsInfo, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalln(err)
	}
	// 根据不同的回调，生成 intents
	intent := event.RegisterHandlers(handler.ATMessageEventHandler(api))
	// 指定需要启动的分片数为2
	wsInfo.Shards = 2
	if err = botgo.NewSessionManager().Start(wsInfo, botToken, &intent); err != nil {
		log.Fatalln(err)
	}
}
