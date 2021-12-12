package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

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
	// 指定需要启动的分片数为2
	wsInfo.Shards = 2
	if err = botgo.NewSessionManager().Start(wsInfo, botToken, &intent); err != nil {
		log.Fatalln(err)
	}
}

// ATMessageEventHandler 实现处理 at 消息的回调
func ATMessageEventHandler(api openapi.OpenAPI) websocket.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		log.Printf("[%s] guildID is %s, content is %s", event.Type, data.GuildID, data.Content)
		if _, err := api.PostMessage(context.Background(), data.ChannelID,
			&dto.MessageToCreate{
				Content: message.MentionAllUser() + fmt.Sprintf("guildID is %s", data.GuildID),
			},
		); err != nil {
			log.Fatalln(err)
		}
		return nil
	}
}
