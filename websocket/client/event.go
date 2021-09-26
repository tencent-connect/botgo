package client

import (
	"encoding/json"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/websocket"
	"github.com/tidwall/gjson" // 由于回包的 d 类型不确定，gjson 用于从回包json中提取 d 并进行针对性的解析
)

var eventParseFuncMap = map[dto.OPCode]map[dto.EventType]eventParseFunc{
	dto.WSDispatchEvent: {
		dto.EventGuildCreate: guildHandler,
		dto.EventGuildUpdate: guildHandler,
		dto.EventGuildDelete: guildHandler,

		dto.EventChannelCreate: channelHandler,
		dto.EventChannelUpdate: channelHandler,
		dto.EventChannelDelete: channelHandler,

		dto.EventGuildMemberAdd:    guildMemberHandler,
		dto.EventGuildMemberUpdate: guildMemberHandler,
		dto.EventGuildMemberRemove: guildMemberHandler,

		dto.EventMessageCreate:   messageHandler,
		dto.EventAtMessageCreate: atMessageHandler,
	},
}

type eventParseFunc func(event *dto.WSPayload, message []byte) error

func parseAndHandle(event *dto.WSPayload, message []byte) error {
	// 指定类型的 handler
	if h, ok := eventParseFuncMap[event.OPCode][event.Type]; ok {
		return h(event, message)
	}
	// 透传handler，如果未注册具体类型的 handler，会统一投递到这个 handler
	if websocket.DefaultHandler.Plain != nil {
		return websocket.DefaultHandler.Plain(event, message)
	}
	return nil
}

func guildHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSGuildData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandler.Guild != nil {
		return websocket.DefaultHandler.Guild(event, data)
	}
	return nil
}

func channelHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSChannelData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandler.Channel != nil {
		return websocket.DefaultHandler.Channel(event, data)
	}
	return nil
}

func guildMemberHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSGuildMemberData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandler.GuildMember != nil {
		return websocket.DefaultHandler.GuildMember(event, data)
	}
	return nil
}

func messageHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSMessageData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandler.Message != nil {
		return websocket.DefaultHandler.Message(event, data)
	}
	return nil
}

func atMessageHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSATMessageData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandler.Message != nil {
		return websocket.DefaultHandler.ATMessage(event, data)
	}
	return nil
}

func parseData(message []byte, target interface{}) error {
	data := gjson.Get(string(message), "d")
	return json.Unmarshal([]byte(data.String()), target)
}
