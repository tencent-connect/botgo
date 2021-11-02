package testcase

import (
	"log"
	"testing"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/websocket"
)

func Test_websocket(t *testing.T) {
	ws, err := api.WS(ctx, nil, "")
	log.Printf("%+v, err:%v", ws, err)

	var message websocket.MessageEventHandler = func(event *dto.WSPayload, data *dto.WSMessageData) error {
		log.Println(event, data)
		return nil
	}
	intent := websocket.RegisterHandlers(message)
	botgo.NewSessionManager().Start(ws, botToken, &intent)
}
