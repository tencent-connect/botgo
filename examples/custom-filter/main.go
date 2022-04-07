package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
)

func main() {
	ctx := context.Background()
	openapi.RegisterReqFilter("set-trace", ReqFilter)
	openapi.RegisterRespFilter("get-trace", RespFilter)
	// 加载 appid 和 token
	botToken := token.New(token.TypeBot)
	if err := botToken.LoadFromConfig("config.yaml"); err != nil {
		log.Fatalln(err)
	}
	// 初始化 openapi，使用 NewSandboxOpenAPI 请求到沙箱环境
	api := botgo.NewSandboxOpenAPI(botToken).WithTimeout(3 * time.Second)
	// 获取 websocket 信息，如果 api 是请求到沙箱环境的，则获取到沙箱环境的 ws 地址
	// websocket 的链接，以及事件处理，请参考其他 examples
	wsInfo, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(wsInfo.URL)
}

// ReqFilter 自定义请求过滤器
func ReqFilter(req *http.Request, _ *http.Response) error {
	req.Header.Set("X-Custom-TraceID", uuid.NewString())
	return nil
}

// RespFilter 自定义响应过滤器
func RespFilter(req *http.Request, resp *http.Response) error {
	log.Println("trace id added by req filter", req.Header.Get("X-Custom-TraceID"))
	log.Println("trace id return by openapi", resp.Header.Get(openapi.TraceIDKey))
	return nil
}
