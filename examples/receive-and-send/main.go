package main

import (
	"context"
	"fmt"
	"log"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

// 消息处理器，持有 openapi 对象
var processor Processor

func main() {
	ctx := context.Background()
	// 加载 appid 和 token
	botToken := token.New(token.TypeBot)
	if err := botToken.LoadFromConfig(getConfigPath("config.yaml")); err != nil {
		log.Fatalln(err)
	}
	// 初始化 openapi
	api := botgo.NewOpenAPI(botToken).WithTimeout(3 * time.Second)
	// 获取 websocket 信息
	wsInfo, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalln(err)
	}

	processor = Processor{api: api}

	// 根据不同的回调，生成 intents
	intent := websocket.RegisterHandlers(
		ATMessageEventHandler(), ReadyHandler(), ErrorNotifyHandler(),
	)
	if err = botgo.NewSessionManager().Start(wsInfo, botToken, &intent); err != nil {
		log.Fatalln(err)
	}
}

// ReadyHandler 自定义 ReadyHandler 感知连接成功事件
func ReadyHandler() websocket.ReadyHandler {
	return func(event *dto.WSPayload, data *dto.WSReadyData) {
		log.Println("ready event receive: ", data)
	}
}

func ErrorNotifyHandler() websocket.ErrorNotifyHandler {
	return func(err error) {
		log.Println("error notify receive: ", err)
	}
}

// ATMessageEventHandler 实现处理 at 消息的回调
func ATMessageEventHandler() websocket.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		input := strings.ToLower(message.ETLInput(data.Content))
		return processor.ProcessMessage(input, data)
	}
}

func getConfigPath(name string) string {
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		return fmt.Sprintf("%s/%s", path.Dir(filename), name)
	}
	return ""
}
