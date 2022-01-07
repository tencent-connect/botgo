package websocket

import (
	"github.com/tencent-connect/botgo/dto"
)

// DefaultHandlers 默认的 handler 结构，管理所有支持的 handler 类型
var DefaultHandlers struct {
	Plain           PlainEventHandler
	Guild           GuildEventHandler
	GuildMember     GuildMemberEventHandler
	Channel         ChannelEventHandler
	Message         MessageEventHandler
	MessageReaction MessageReactionEventHandler
	ATMessage       ATMessageEventHandler
	DirectMessage   DirectMessageEventHandler
	Audio           AudioEventHandler
	MessageAudit    MessageAuditEventHandler
}

// PlainEventHandler 透传handler
type PlainEventHandler func(event *dto.WSPayload, message []byte) error

// GuildEventHandler 频道事件handler
type GuildEventHandler func(event *dto.WSPayload, data *dto.WSGuildData) error

// GuildMemberEventHandler 频道成员事件 handler
type GuildMemberEventHandler func(event *dto.WSPayload, data *dto.WSGuildMemberData) error

// ChannelEventHandler 子频道事件 handler
type ChannelEventHandler func(event *dto.WSPayload, data *dto.WSChannelData) error

// MessageEventHandler 消息事件 handler
type MessageEventHandler func(event *dto.WSPayload, data *dto.WSMessageData) error

// MessageReactionEventHandler 表情表态事件 handler
type MessageReactionEventHandler func(event *dto.WSPayload, data *dto.WSMessageReactionData) error

// ATMessageEventHandler at 机器人消息事件 handler
type ATMessageEventHandler func(event *dto.WSPayload, data *dto.WSATMessageData) error

// DirectMessageEventHandler 私信消息事件 handler
type DirectMessageEventHandler func(event *dto.WSPayload, data *dto.WSDirectMessageData) error

// AudioEventHandler 音频机器人事件 handler
type AudioEventHandler func(event *dto.WSPayload, data *dto.WSAudioData) error

// MessageAuditEventHandler 消息审核事件 handler
type MessageAuditEventHandler func(event *dto.WSPayload, data *dto.WSMessageAuditData) error

// RegisterHandlers 注册事件回调，并返回 intent 用于 websocket 的鉴权
func RegisterHandlers(handlers ...interface{}) dto.Intent {
	var i dto.Intent
	for _, h := range handlers {
		switch handle := h.(type) {
		case PlainEventHandler:
			DefaultHandlers.Plain = handle
		case AudioEventHandler:
			DefaultHandlers.Audio = handle
			i = i | dto.EventToIntent(
				dto.EventAudioStart, dto.EventAudioFinish,
				dto.EventAudioOnMic, dto.EventAudioOffMic,
			)
		default:
		}
	}
	i = i | registerGuildHandlers(i, handlers...)
	i = i | registerMessageHandlers(i, handlers...)

	return i
}

// registerGuildHandlers 注册频道相关handlers
func registerGuildHandlers(i dto.Intent, handlers ...interface{}) dto.Intent {
	for _, h := range handlers {
		switch handle := h.(type) {
		case GuildEventHandler:
			DefaultHandlers.Guild = handle
			i = i | dto.EventToIntent(dto.EventGuildCreate, dto.EventGuildDelete, dto.EventGuildUpdate)
		case GuildMemberEventHandler:
			DefaultHandlers.GuildMember = handle
			i = i | dto.EventToIntent(dto.EventGuildMemberAdd, dto.EventGuildMemberRemove, dto.EventGuildMemberUpdate)
		case ChannelEventHandler:
			DefaultHandlers.Channel = handle
			i = i | dto.EventToIntent(dto.EventChannelCreate, dto.EventChannelDelete, dto.EventChannelUpdate)
		default:
		}
	}
	return i
}

// registerMessageHandlers 注册消息相关的 handler
func registerMessageHandlers(i dto.Intent, handlers ...interface{}) dto.Intent {
	for _, h := range handlers {
		switch handle := h.(type) {
		case MessageEventHandler:
			DefaultHandlers.Message = handle
			i = i | dto.EventToIntent(dto.EventMessageCreate)
		case ATMessageEventHandler:
			DefaultHandlers.ATMessage = handle
			i = i | dto.EventToIntent(dto.EventAtMessageCreate)
		case DirectMessageEventHandler:
			DefaultHandlers.DirectMessage = handle
			i = i | dto.EventToIntent(dto.EventDirectMessageCreate)
		case MessageReactionEventHandler:
			DefaultHandlers.MessageReaction = handle
			i = i | dto.EventToIntent(dto.EventMessageReactionAdd, dto.EventMessageReactionRemove)
		case MessageAuditEventHandler:
			DefaultHandlers.MessageAudit = handle
			i = i | dto.EventToIntent(dto.EventMessageAuditPass, dto.EventMessageAuditReject)
		default:
		}
	}
	return i
}
