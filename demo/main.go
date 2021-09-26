package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	"gopkg.in/yaml.v2"
)

var conf struct {
	AppID uint64 `yaml:"appid"`
	Token string `yaml:"token"`
}

func init() {
	content, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Println("read conf failed")
		os.Exit(1)
	}
	if err := yaml.Unmarshal(content, &conf); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println(conf)
}

type runFunc func() (interface{}, error)

func main() {
	token := token.BotToken(conf.AppID, conf.Token)
	api := botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
	ctx := context.Background()
	ws, err := api.WS(ctx, nil, "")
	log.Printf("%+v, err:%v", ws, err)

	// me, err := api.Me(ctx, nil, "")
	// log.Printf("%+v, err:%v", me, err)
	//
	// run(func() (interface{}, error) {
	// 	return api.MeGuilds(ctx)
	// })

	// run(func() (interface{}, error) {
	// 	return api.Guild(ctx, "13034202056525133443")
	// })

	// run(func() (interface{}, error) {
	// 	return api.GuildMembers(ctx, "13034202056525133443", &dto.GuildMembersPager{
	// 		After: "0",
	// 		Limit: "100",
	// 	})
	// })

	// run(func() (interface{}, error) {
	// 	return api.Channels(ctx, "13034202056525133443")
	// })
	//
	// run(func() (interface{}, error) {
	// 	return api.Channel(ctx, "1107217")
	// })
	//
	// run(func() (interface{}, error) {
	// 	return api.PostMessage(ctx, "1107217", &dto.MessageToCreate{Content: "abc"})
	// })

	var message websocket.MessageEventHandler = func(event *dto.WSPayload, data *dto.WSMessageData) error {
		log.Println(event, data)
		return nil
	}
	intent := websocket.RegisterHandlers(message)
	botgo.NewSessionManager().Start(ws, token, &intent)
}

func run(runFunc2 runFunc) {
	ret, err := runFunc2()
	log.Printf("%+v, err:%v", ret, err)
}
