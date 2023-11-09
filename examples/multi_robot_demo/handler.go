package main

import (
	"context"
	"fmt"
	"multi_robot/processor"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
)

// InteractionEventHandler 互动事件按钮点击
func InteractionEventHandler() event.InteractionEventHandler {
	return func(event *dto.WSPayload, data *dto.WSInteractionData) error {
		fmt.Printf("%+v\n", data)
		return nil
	}
}

// C2CMessageEventHandler 实现处理 at 消息的回调
func C2CMessageEventHandler() event.C2CMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSC2CMessageData) error {
		fmt.Println("%+v\n", data)
		userID := ""
		if data.Author != nil && data.Author.ID != "" {
			userID = data.Author.ID
		}

		proc := processor.GetProcessor(event.Session.Token.GetAppID())
		msg, _ := proc.GetReplayMsg(context.Background(), dto.Message(*data))

		proc.GetAPI().PostC2CMessage(context.Background(), userID, msg)
		return nil
	}
}

// GroupATMessageEventHandler 实现处理 at 消息的回调
func GroupATMessageEventHandler() event.GroupATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		fmt.Printf("%+v\n", data)

		proc := processor.GetProcessor(event.Session.Token.GetAppID())
		msg, _ := proc.GetReplayMsg(context.Background(), dto.Message(*data))
		proc.
			GetAPI().PostGroupMessage(context.Background(), data.GroupID, msg)
		return nil

	}
}
