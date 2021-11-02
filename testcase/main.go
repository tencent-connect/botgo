package testcase

import (
	"context"
	"log"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
)

var conf struct {
	AppID uint64 `yaml:"appid"`
	Token string `yaml:"token"`
}
var botToken *token.Token
var api openapi.OpenAPI

type runFunc func() (interface{}, error)

func main() {
	token := token.BotToken(conf.AppID, conf.Token)
	api := botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
	ctx := context.Background()

	me, err := api.Me(ctx)
	log.Printf("%+v, err:%v", me, err)
	//
	// run(func() (interface{}, error) {
	// 	return api.MeGuilds(ctx)
	// })

	// run(func() (interface{}, error) {
	// 	return api.Guild(ctx, "3326534247441079828")
	// })

	// run(func() (interface{}, error) {
	// 	return api.Roles(ctx, "3326534247441079828")
	// })

	// run(func() (interface{}, error) {
	// 	return api.PostRole(ctx, "3326534247441079828", &dto.Role{
	// 		Name:        "test roles",
	// 		Color:       4278245297,
	// 		Hoist:       1,
	// 		MemberCount: 1,
	// 		MemberLimit: 1,
	// 	})
	// })

	// run(func() (interface{}, error) {
	// 	return api.PatchRole(ctx, "3326534247441079828", "10002894", &dto.Role{
	// 		Name:        "test roles 22",
	// 		Color:       4278245297,
	// 		Hoist:       1,
	// 		MemberCount: 1,
	// 		MemberLimit: 1,
	// 	})
	// })
	//
	// run(func() (interface{}, error) {
	// 	return nil, api.DeleteRole(ctx, "3326534247441079828", "10002894")
	// })

	// run(func() (interface{}, error) {
	// 	return api.GuildMembers(ctx, "3326534247441079828", &dto.GuildMembersPager{
	// 		After: "0",
	// 		Limit: "1",
	// 	})
	// })

	// run(func() (interface{}, error) {
	// 	return nil, api.MemberAddRole(ctx, "3326534247441079828", "10002878", "6103844827365541271")
	// })
	// run(func() (interface{}, error) {
	// 	return nil, api.MemberDeleteRole(ctx, "3326534247441079828", "10002878", "6103844827365541271")
	// })

	// run(func() (interface{}, error) {
	// 	return api.Channels(ctx, "3326534247441079828")
	// })
	//
	// run(func() (interface{}, error) {
	// 	return api.Channel(ctx, "1107217")
	// })
	//
	// run(func() (interface{}, error) {
	// 	return api.PostMessage(ctx, "1107217", &dto.MessageToCreate{Content: "abc"})
	// })
}

func run(runFunc2 runFunc) {
	ret, err := runFunc2()
	log.Printf("%+v, err:%v", ret, err)
}
