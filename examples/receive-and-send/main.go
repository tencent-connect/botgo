package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
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

	// 初始化 openapi，正式环境
	api := botgo.NewOpenAPI(botToken).WithTimeout(3 * time.Second)
	// 沙箱环境
	// api := botgo.NewSandboxOpenAPI(botToken).WithTimeout(3 * time.Second)

	// 获取 websocket 信息
	wsInfo, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalln(err)
	}

	processor = Processor{api: api}

	//websocket.RegisterResumeSignal(syscall.SIGUSR1)
	// 根据不同的回调，生成 intents
	intent := websocket.RegisterHandlers(
		// at 机器人事件，目前是在这个事件处理中有逻辑，会回消息，其他的回调处理都只把数据打印出来，不做任何处理
		ATMessageEventHandler(),
		// 如果想要捕获到连接成功的事件，可以实现这个回调
		ReadyHandler(),
		// 连接关闭回调
		ErrorNotifyHandler(),
		// 频道事件
		GuildEventHandler(),
		// 成员事件
		MemberEventHandler(),
		// 子频道事件
		ChannelEventHandler(),
		// 私信，目前只有私域才能够收到这个，如果你的机器人不是私域机器人，会导致连接报错，那么启动 example 就需要注释掉这个回调
		DirectMessageHandler(),
		// 频道消息，只有私域才能够收到这个，如果你的机器人不是私域机器人，会导致连接报错，那么启动 example 就需要注释掉这个回调
		CreateMessageHandler(),
		// 互动事件
		InteractionHandler(api),
	)
	// 指定需要启动的分片数为 2 的话可以手动修改 wsInfo
	// wsInfo.Shards = 2
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

func GuildEventHandler() websocket.GuildEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGuildData) error {
		fmt.Println(data)
		return nil
	}
}

func ChannelEventHandler() websocket.ChannelEventHandler {
	return func(event *dto.WSPayload, data *dto.WSChannelData) error {
		fmt.Println(data)
		return nil
	}
}

func MemberEventHandler() websocket.GuildMemberEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGuildMemberData) error {
		fmt.Println(data)
		return nil
	}
}

func DirectMessageHandler() websocket.DirectMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSDirectMessageData) error {
		fmt.Println(data)
		return nil
	}
}

func CreateMessageHandler() websocket.MessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSMessageData) error {
		fmt.Println(data)
		return nil
	}
}

func InteractionHandler(api openapi.OpenAPI) websocket.InteractionEventHandler {
	return func(event *dto.WSPayload, data *dto.WSInteractionData) error {
		fmt.Println(data)
		if err := putInteractionInlineSearch(api, data); err != nil {
			log.Println("putInteractionInlineSearch error: ", err)
		}
		return nil
	}
}

func getConfigPath(name string) string {
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		return fmt.Sprintf("%s/%s", path.Dir(filename), name)
	}
	return ""
}

func putInteractionInlineSearch(api openapi.OpenAPI, interaction *dto.WSInteractionData) error {
	if interaction.Data.Type != dto.InteractionDataTypeChatSearch {
		return errors.New("interaction data type not chat search")
	}
	search := &dto.SearchInputResolved{}
	if err := json.Unmarshal(interaction.Data.Resolved, search); err != nil {
		return fmt.Errorf("interaction data resolved unmarshal error:%v", err)
	}
	if search.Keyword != "test" {
		return errors.New("resolved search key not allowed")
	}
	searchRsp := &dto.SearchRsp{
		Layouts: []dto.SearchLayout{
			{
				LayoutType: 0,
				ActionType: 0,
				Title:      "内联搜索",
				Records: []dto.SearchRecord{
					{
						Cover: "https://pub.idqqimg.com/pc/misc/files/20211208/311cfc87ce394c62b7c9f0508658cf25.png",
						Title: "内联搜索标题",
						Tips:  "内联搜索 tips",
						URL:   "https://www.qq.com",
					},
				},
			},
		},
	}
	body, _ := json.Marshal(searchRsp)
	if err := api.PutInteraction(context.Background(), interaction.ID, body); err != nil {
		log.Println("api call putInteractionInlineSearch  error: ", err)
		return err
	}
	return nil
}
