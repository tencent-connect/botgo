package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
)

func main() {
	ctx := context.Background()

	// 初始化新的文件 logger，并使用相对路径来作为日志存放位置，设置最小日志界别为 DebugLevel
	logger, err := New("./", DebugLevel)
	if err != nil {
		log.Fatalln("error log new", err)
	}
	// 把新的 logger 设置到 sdk 上，替换掉老的控制台 logger
	botgo.SetLogger(logger)

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
	intent := event.RegisterHandlers(ATMessageEventHandler(api))
	if err = botgo.NewSessionManager().Start(wsInfo, botToken, &intent); err != nil {
		log.Fatalln(err)
	}
}

// ATMessageEventHandler 实现处理 at 消息的回调
func ATMessageEventHandler(api openapi.OpenAPI) event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		log.Printf("[%s] %s", event.Type, data.Content)
		input := strings.ToLower(message.ETLInput(data.Content))
		log.Printf("clear input content is: %s", input)
		return nil
	}
}
