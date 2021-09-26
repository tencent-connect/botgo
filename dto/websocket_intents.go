package dto

// Intent 类型
type Intent int

// websocket intent 声明
const (
	// IntentGuilds 包含
	// - GUILD_CREATE
	// - GUILD_UPDATE
	// - GUILD_DELETE
	// - GUILD_ROLE_CREATE
	// - GUILD_ROLE_UPDATE
	// - GUILD_ROLE_DELETE
	// - CHANNEL_CREATE
	// - CHANNEL_UPDATE
	// - CHANNEL_DELETE
	// - CHANNEL_PINS_UPDATE
	IntentGuilds Intent = 1 << iota

	// IntentGuildMembers 包含
	// - GUILD_MEMBER_ADD
	// - GUILD_MEMBER_UPDATE
	// - GUILD_MEMBER_REMOVE
	IntentGuildMembers

	IntentGuildBans
	IntentGuildEmojis
	IntentGuildIntegrations
	IntentGuildWebhooks
	IntentGuildInvites
	IntentGuildVoiceStates
	IntentGuildPresences
	IntentGuildMessages
	IntentGuildMessageReactions
	IntentGuildMessageTyping
	IntentDirectMessages
	IntentDirectMessageReactions
	IntentDirectMessageTyping

	IntentsGuildAtMessage Intent = 1 << 30 // 只接收@消息事件

	IntentNone Intent = 0
)

var eventIntentMap = map[EventType]Intent{
	EventGuildCreate:   IntentGuilds,
	EventGuildUpdate:   IntentGuilds,
	EventGuildDelete:   IntentGuilds,
	EventChannelCreate: IntentGuilds,
	EventChannelUpdate: IntentGuilds,
	EventChannelDelete: IntentGuilds,

	EventGuildMemberAdd:    IntentGuildMembers,
	EventGuildMemberUpdate: IntentGuildMembers,
	EventGuildMemberRemove: IntentGuildMembers,

	EventMessageCreate:   IntentGuildMessages,
	EventAtMessageCreate: IntentsGuildAtMessage,
}

// EventToIntent 事件转换对应的Intent
func EventToIntent(events ...EventType) Intent {
	var i Intent
	for _, event := range events {
		i = i | eventIntentMap[event]
	}
	return i
}
