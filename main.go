package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

func main() {
	token := token.BotToken(0, "", string(token.TypeBot))
	api := NewOpenAPI(token).WithTimeout(3 * time.Second)
	ctx := context.Background()
	ws, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Warn("登录失败，请检查 appid 和 AccessToken 是否正确。")
		log.Info("该程序将于5秒后退出！")
		time.Sleep(time.Second * 5)
	}
	var groupMessage event.GroupAtMessageEventHandler = func(event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		resp, _ := api.PostGroupRichMediaMessage(ctx, data.GroupId, &dto.GroupRichMediaMessageToCreate{FileType: 1, Url: "https://www.2mf8.cn/static/image/cube3/b1.png", SrvSendMsg: false})
		if resp != nil {
			newMsg := &dto.GroupMessageToCreate{
				Content: "msg", //+ "\n[🔗奇乐最新价格]\n(https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F)",
				Media: &dto.FileInfo{
					FileInfo: resp.FileInfo,
				},
				MsgID:   data.MsgId,
				MsgType: 7,
				MsgReq:  1,
			}
			api.PostGroupMessage(ctx, data.GroupId, newMsg)
		}
		return nil
	}
	intent := websocket.RegisterHandlers(groupMessage)
	NewSessionManager().Start(ws, token, &intent)
}
