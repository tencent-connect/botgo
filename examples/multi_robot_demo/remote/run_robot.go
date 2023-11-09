// Package remote 多机器人使用demo
package remote

import (
	"context"
	"fmt"
	"multi_robot/processor"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/sessions/remote"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	"gopkg.in/yaml.v2"
)

type addProcessor func(uint64, processor.Processor)
type processorBuilder func(uint64) processor.Processor

type robotDetail struct {
	AppID  uint64 `yaml:"appid"`
	AppKey string `yaml:"appkey"`
	Redis  struct {
		PassWord string `yaml:"pass_word"`
		User     string `yaml:"user"`
		Addr     string `yaml:"addr"`
		NetWork  string `yaml:"net_work"`
	} `yaml:"redis"`

	redisInstance *redis.Client
	ClusterKey    string `yaml:"cluster_key"`
}

// 机器人列表读取器
func getRobotList(configFile string) []*robotDetail {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	var cfg []*robotDetail
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		panic(err)
	}

	for _, f := range cfg {
		f.redisInstance = redis.NewClient(&redis.Options{
			Password: f.Redis.PassWord,
			Username: f.Redis.User,
			Network:  f.Redis.NetWork,
			Addr:     f.Redis.Addr,
		})
	}

	fmt.Printf("%+v\n", cfg)
	return cfg
}

// InitProcessRobot 启动机器人
func InitProcessRobot(configFile string,
	add addProcessor, pb processorBuilder, handlers ...interface{}) {
	ctx := context.Background()

	for _, v := range getRobotList(configFile) {
		robotInfo := v
		go initNewRobotProcess(ctx, robotInfo, add, pb, handlers...)
	}
}

// InitNewRobotProcess 启动单个机器人
func initNewRobotProcess(ctx context.Context, robotDetail *robotDetail,
	add addProcessor, pb processorBuilder,
	handlers ...interface{}) {
	// 加载 appid 和 token

	botToken := token.BotToken(robotDetail.AppID, robotDetail.AppKey)
	botToken.Type = token.TypeQQBot

	botgo.NewSessionManager()
	if err := botToken.InitToken(ctx); err != nil {
		fmt.Printf("%v,robot:%v\n", err, robotDetail)
		return
	}
	//初始化 openapi，正式环境
	api := botgo.NewOpenAPI(botToken).WithTimeout(3 * time.Second)
	// 获取 websocket 信息
	wsInfo, err := api.WS(ctx, nil, "")
	if err != nil {
		fmt.Printf("init robot appid[%d] invalid,err:%v\n", robotDetail.AppID, err)
		return
	}

	proc := pb(robotDetail.AppID)
	proc.SetAPI(api)
	add(robotDetail.AppID, proc)
	// 根据不同的回调，生成 intents
	intent := websocket.RegisterHandlers(
		handlers...,
	)
	// 指定需要启动的分片数为 2 的话可以手动修改 wsInfo
	if err = remote.New(robotDetail.redisInstance,
		remote.WithClusterKey(robotDetail.ClusterKey)).Start(wsInfo, botToken, &intent); err != nil {
		fmt.Printf("%v,robot:%v\n", err, robotDetail)
	}
}
