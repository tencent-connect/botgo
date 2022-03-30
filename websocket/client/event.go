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

		dto.EventMessageCreate: messageHandler,

		dto.EventMessageReactionAdd:    messageReactionHandler,
		dto.EventMessageReactionRemove: messageReactionHandler,

		dto.EventAtMessageCreate:     atMessageHandler,
		dto.EventDirectMessageCreate: directMessageHandler,

		dto.EventAudioStart:  audioHandler,
		dto.EventAudioFinish: audioHandler,
		dto.EventAudioOnMic:  audioHandler,
		dto.EventAudioOffMic: audioHandler,

		dto.EventMessageAuditPass:   messageAuditHandler,
		dto.EventMessageAuditReject: messageAuditHandler,

		dto.EventForumThreadCreate: threadHandler,
		dto.EventForumThreadUpdate: threadHandler,
		dto.EventForumThreadDelete: threadHandler,
		dto.EventForumPostCreate:   postHandler,
		dto.EventForumPostDelete:   postHandler,
		dto.EventForumReplyCreate:  replyHandler,
		dto.EventForumReplyDelete:  replyHandler,
		dto.EventForumAuditResult:  forumAuditHandler,
		dto.EventInteractionCreate: interactionHandler,
	},
}

type eventParseFunc func(event *dto.WSPayload, message []byte) error

func parseAndHandle(event *dto.WSPayload) error {
	// 指定类型的 handler
	if h, ok := eventParseFuncMap[event.OPCode][event.Type]; ok {
		return h(event, event.RawMessage)
	}
	// 透传handler，如果未注册具体类型的 handler，会统一投递到这个 handler
	if websocket.DefaultHandlers.Plain != nil {
		return websocket.DefaultHandlers.Plain(event, event.RawMessage)
	}
	return nil
}

func guildHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSGuildData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandlers.Guild != nil {
		return websocket.DefaultHandlers.Guild(event, data)
	}
	return nil
}

func channelHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSChannelData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandlers.Channel != nil {
		return websocket.DefaultHandlers.Channel(event, data)
	}
	return nil
}

func guildMemberHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSGuildMemberData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandlers.GuildMember != nil {
		return websocket.DefaultHandlers.GuildMember(event, data)
	}
	return nil
}

func messageHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSMessageData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandlers.Message != nil {
		return websocket.DefaultHandlers.Message(event, data)
	}
	return nil
}

func messageReactionHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSMessageReactionData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandlers.MessageReaction != nil {
		return websocket.DefaultHandlers.MessageReaction(event, data)
	}
	return nil
}

func atMessageHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSATMessageData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandlers.ATMessage != nil {
		return websocket.DefaultHandlers.ATMessage(event, data)
	}
	return nil
}

func directMessageHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSDirectMessageData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandlers.DirectMessage != nil {
		return websocket.DefaultHandlers.DirectMessage(event, data)
	}
	return nil
}

func audioHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSAudioData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandlers.Audio != nil {
		return websocket.DefaultHandlers.Audio(event, data)
	}
	return nil
}

func threadHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSThreadData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandlers.Thread != nil {
		return websocket.DefaultHandlers.Thread(event, data)
	}
	return nil
}

func postHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSPostData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandlers.Post != nil {
		return websocket.DefaultHandlers.Post(event, data)
	}
	return nil
}

func replyHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSReplyData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandlers.Reply != nil {
		return websocket.DefaultHandlers.Reply(event, data)
	}
	return nil
}

func forumAuditHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSForumAuditData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandlers.ForumAudit != nil {
		return websocket.DefaultHandlers.ForumAudit(event, data)
	}
	return nil
}

func interactionHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSInteractionData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandlers.Interaction != nil {
		return websocket.DefaultHandlers.Interaction(event, data)
	}
	return nil
}

func parseData(message []byte, target interface{}) error {
	data := gjson.Get(string(message), "d")
	return json.Unmarshal([]byte(data.String()), target)
}

func messageAuditHandler(event *dto.WSPayload, message []byte) error {
	data := &dto.WSMessageAuditData{}
	if err := parseData(message, data); err != nil {
		return err
	}
	if websocket.DefaultHandlers.MessageAudit != nil {
		return websocket.DefaultHandlers.MessageAudit(event, data)
	}
	return nil
}
