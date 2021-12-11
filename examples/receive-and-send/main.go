package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

// 输入输出词典
var dict = map[string]string{
	"hello": "World",
	"hi":    "NoHi",
}

func main() {
	ctx := context.Background()
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
	intent := websocket.RegisterHandlers(ATMessageEventHandler(api))
	if err = botgo.NewSessionManager().Start(wsInfo, botToken, &intent); err != nil {
		log.Fatalln(err)
	}
}

// ATMessageEventHandler 实现处理 at 消息的回调
func ATMessageEventHandler(api openapi.OpenAPI) websocket.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		log.Printf("[%s] %s", event.Type, data.Content)
		input := strings.ToLower(message.ETLInput(data.Content))
		log.Printf("clear input content is: %s", input)
		// 根据词典中的输入，进行输出
		if v, ok := dict[input]; ok {
			if _, err := api.PostMessage(context.Background(), data.ChannelID,
				&dto.MessageToCreate{
					Content: message.MentionUser(data.Author.ID) + v,
				},
			); err != nil {
				log.Fatalln(err)
			}
		}
		return nil
	}
}
