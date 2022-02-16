package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
)

type Processor struct {
	api openapi.OpenAPI
}

func (p Processor) ProcessMessage(input string, data *dto.WSATMessageData) error {
	cmd := message.ParseCommand(input)
	toCreate := &dto.MessageToCreate{
		Content: "默认回复 <emoji:37>",
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
	}
	// 进入到私信逻辑
	if cmd.Cmd == "dm" {
		p.dmHandler(data)
		return nil
	}

	var reply bool
	switch cmd.Cmd {
	case "time":
		toCreate.Content = genReplyContent(data)
		reply = true
	case "ark":
		toCreate.Ark = genReplyArk(data)
		reply = true
	case "hi":
		reply = true
	default:
	}
	fmt.Println("1111111", reply, []byte(cmd.Cmd), len(cmd.Cmd))
	// 是否命中上面的指令，不回复多余的内容
	if reply {
		if _, err := p.api.PostMessage(context.Background(), data.ChannelID, toCreate); err != nil {
			log.Println(err)
		}
	}

	return nil
}

func (p Processor) dmHandler(data *dto.WSATMessageData) {
	dm, err := p.api.CreateDirectMessage(
		context.Background(), &dto.DirectMessageToCreate{
			SourceGuildID: data.GuildID,
			RecipientID:   data.Author.ID,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}

	toCreate := &dto.MessageToCreate{
		Content: "默认私信回复",
	}
	_, err = p.api.PostDirectMessage(
		context.Background(), dm, toCreate,
	)
	if err != nil {
		log.Println(err)
		return
	}
}

func genReplyContent(data *dto.WSATMessageData) string {
	var tpl = `你好：%s
在子频道 %s 收到消息。
收到的消息发送时时间为：%s
当前本地时间为：%s

消息来自：%s
`

	msgTime, _ := data.Timestamp.Time()
	return fmt.Sprintf(
		tpl,
		message.MentionUser(data.Author.ID),
		message.MentionChannel(data.ChannelID),
		msgTime, time.Now().Format(time.RFC3339),
		getIP(),
	)
}

func genReplyArk(data *dto.WSATMessageData) *dto.Ark {
	return &dto.Ark{
		TemplateID: 23,
		KV: []*dto.ArkKV{
			{
				Key:   "#DESC#",
				Value: "这是 ark 的描述信息",
			},
			{
				Key:   "#PROMPT#",
				Value: "这是 ark 的摘要信息",
			},
			{
				Key: "#LIST#",
				Obj: []*dto.ArkObj{
					{
						ObjKV: []*dto.ArkObjKV{
							{
								Key:   "desc",
								Value: "这里展示的是 23 号模板",
							},
						},
					},
					{
						ObjKV: []*dto.ArkObjKV{
							{
								Key:   "desc",
								Value: "这是 ark 的列表项名称",
							},
							{
								Key:   "link",
								Value: "https://www.qq.com",
							},
						},
					},
				},
			},
		},
	}
}
