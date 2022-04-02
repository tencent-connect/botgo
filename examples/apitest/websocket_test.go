package apitest

import (
	"log"
	"testing"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
)

func Test_websocket(t *testing.T) {
	ws, err := api.WS(ctx, nil, "")
	log.Printf("%+v, err:%v", ws, err)

	t.Run(
		"at message", func(t *testing.T) {
			var message event.ATMessageEventHandler = func(event *dto.WSPayload, data *dto.WSATMessageData) error {
				log.Println(event, data)
				return nil
			}
			intent := event.RegisterHandlers(message)
			botgo.NewSessionManager().Start(ws, botToken, &intent)
		},
	)
	t.Run(
		"at message assign shard to 2", func(t *testing.T) {
			var message event.ATMessageEventHandler = func(event *dto.WSPayload, data *dto.WSATMessageData) error {
				log.Println(event, data)
				return nil
			}
			ws.Shards = 2
			intent := event.RegisterHandlers(message)
			botgo.NewSessionManager().Start(ws, botToken, &intent)
		},
	)
	t.Run(
		"at message and guild event", func(t *testing.T) {
			var message event.ATMessageEventHandler = func(event *dto.WSPayload, data *dto.WSATMessageData) error {
				log.Println(event, data)
				return nil
			}
			var guildEvent event.GuildEventHandler = func(event *dto.WSPayload, data *dto.WSGuildData) error {
				log.Println(event, data)
				return nil
			}
			intent := event.RegisterHandlers(message, guildEvent)
			botgo.NewSessionManager().Start(ws, botToken, &intent)
		},
	)
	t.Run(
		"message reaction", func(t *testing.T) {
			var message event.MessageReactionEventHandler = func(
				event *dto.WSPayload, data *dto.WSMessageReactionData,
			) error {
				log.Println(event, data)
				return nil
			}
			intent := event.RegisterHandlers(message)
			botgo.NewSessionManager().Start(ws, botToken, &intent)
		},
	)
	t.Run(
		"thread event", func(t *testing.T) {
			var message event.ThreadEventHandler = func(
				event *dto.WSPayload, data *dto.WSThreadData,
			) error {
				log.Println(event, data)
				return nil
			}
			intent := event.RegisterHandlers(message)
			botgo.NewSessionManager().Start(ws, botToken, &intent)
		},
	)
	t.Run(
		"post event", func(t *testing.T) {
			var message event.PostEventHandler = func(
				event *dto.WSPayload, data *dto.WSPostData,
			) error {
				log.Println(event, data)
				return nil
			}
			intent := event.RegisterHandlers(message)
			botgo.NewSessionManager().Start(ws, botToken, &intent)
		},
	)
	t.Run(
		"Reply event", func(t *testing.T) {
			var message event.ReplyEventHandler = func(
				event *dto.WSPayload, data *dto.WSReplyData,
			) error {
				log.Println(event, data)
				return nil
			}
			intent := event.RegisterHandlers(message)
			botgo.NewSessionManager().Start(ws, botToken, &intent)
		},
	)
	t.Run(
		"Forum audit event", func(t *testing.T) {
			var message event.ForumAuditEventHandler = func(
				event *dto.WSPayload, data *dto.WSForumAuditData,
			) error {
				log.Println(event, data)
				return nil
			}
			intent := event.RegisterHandlers(message)
			botgo.NewSessionManager().Start(ws, botToken, &intent)
		},
	)
}
