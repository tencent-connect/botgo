package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/openapi"
)

// ATMessageEventHandler 实现处理 at 消息的回调
func ATMessageEventHandler(api openapi.OpenAPI) event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		log.Printf("[%s] guildID is %s, content is %s", event.Type, data.GuildID, data.Content)
		if _, err := api.PostMessage(
			context.Background(), data.ChannelID,
			&dto.MessageToCreate{
				Content: message.MentionAllUser() + fmt.Sprintf("guildID is %s", data.GuildID),
			},
		); err != nil {
			log.Println(err)
		}
		return nil
	}
}

// GetRedisConn 获取 redis 连接，这只是一个例子，生产中，redis 应该开启 auth，同时 redis 的地址信息应该从配置文件中获取
func GetRedisConn(addr string, timeout time.Duration) *redis.Client {
	return redis.NewClient(
		&redis.Options{
			Addr:         addr,
			DialTimeout:  timeout,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
		},
	)
}
