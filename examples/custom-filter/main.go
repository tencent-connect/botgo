package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/interaction/webhook"
	"gopkg.in/yaml.v3"

	"github.com/tencent-connect/botgo/constant"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
)

const (
	host_ = "0.0.0.0"
	port_ = 9000
	path_ = "/qqbot"
)

func main() {
	ctx := context.Background()
	openapi.RegisterReqFilter("set-trace", ReqFilter)
	openapi.RegisterRespFilter("get-trace", RespFilter)
	// 加载 appid 和 token
	content, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalln("load config file failed, err:", err)
	}
	credentials := &token.QQBotCredentials{}
	if err = yaml.Unmarshal(content, &credentials); err != nil {
		log.Fatalln("parse config failed, err:", err)
	}
	tokenSource := token.NewQQBotTokenSource(credentials)
	if err = token.StartRefreshAccessToken(ctx, tokenSource); err != nil {
		log.Fatalln(err)
	}
	// 初始化 openapi，正式环境
	api := botgo.NewOpenAPI(credentials.AppID, tokenSource).WithTimeout(5 * time.Second).SetDebug(true)
	// 根据不同的回调，生成 intents
	_ = event.RegisterHandlers(GuildATMessageEventHandler(api))
	// 初始化 openapi，正式环境
	http.HandleFunc(path_, func(writer http.ResponseWriter, request *http.Request) {
		webhook.HTTPHandler(writer, request, credentials)
	})
	if err = http.ListenAndServe(fmt.Sprintf("%s:%d", host_, port_), nil); err != nil {
		log.Fatal("setup server fatal:", err)
	}
}

// ReqFilter 自定义请求过滤器
func ReqFilter(req *http.Request, _ *http.Response) error {
	req.Header.Set("X-Custom-TraceID", uuid.NewString())
	return nil
}

// RespFilter 自定义响应过滤器
func RespFilter(req *http.Request, resp *http.Response) error {
	log.Println("trace id added by req filter", req.Header.Get("X-Custom-TraceID"))
	log.Println("trace id return by openapi", resp.Header.Get(constant.HeaderTraceID))
	return nil
}

// GuildATMessageEventHandler 实现处理 at 消息的回调
func GuildATMessageEventHandler(api openapi.OpenAPI) event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		log.Printf("[%s] %s", event.Type, data.Content)
		input := strings.ToLower(message.ETLInput(data.Content))
		log.Printf("clear input content is: %s", input)
		return nil
	}
}
