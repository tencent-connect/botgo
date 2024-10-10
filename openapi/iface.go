package openapi

import (
	"context"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi/options"
	"golang.org/x/oauth2"
)

// OpenAPI openapi 完整实现
type OpenAPI interface {
	Base
	WebsocketAPI
	UserAPI
	MessageAPI
	DirectMessageAPI
	GuildAPI
	ChannelAPI
	AudioAPI
	RoleAPI
	MemberAPI
	ChannelPermissionsAPI
	AnnouncesAPI
	ScheduleAPI
	APIPermissionsAPI
	PinsAPI
	MessageReactionAPI
	WebhookAPI
	InteractionAPI
	MessageSettingAPI
}

// Base 基础能力接口
type Base interface {
	// Version 接口版本
	Version() APIVersion

	// Setup 初始化
	Setup(botAppID string, token oauth2.TokenSource, inSandbox bool) OpenAPI

	// WithTimeout 设置请求接口超时时间
	WithTimeout(duration time.Duration) OpenAPI

	// SetDebug 设置调试模式, 输出更多过程日志
	SetDebug(debug bool) OpenAPI

	// Transport 透传请求，如果 sdk 没有及时跟进新的接口的变更，可以使用该方法进行透传，openapi 实现时可以按需选择是否实现该接口
	Transport(ctx context.Context, method, url string, body interface{}) ([]byte, error)

	// TraceID 返回上一次请求的 trace id
	TraceID() string
}

// WebsocketAPI websocket 接入地址
type WebsocketAPI interface {
	WS(ctx context.Context, params map[string]string, body string) (*dto.WebsocketAP, error)
}

// UserAPI 用户相关接口
type UserAPI interface {
	Me(ctx context.Context) (*dto.User, error)
	MeGuilds(ctx context.Context, pager *dto.GuildPager) ([]*dto.Guild, error)
}

// MessageAPI 消息相关接口
type MessageAPI interface {
	// Message 拉取频道单条消息
	Message(ctx context.Context, channelID string, messageID string, opt ...options.Option) (
		*dto.Message, error)

	// Messages 拉取频道消息列表
	Messages(ctx context.Context, channelID string, pager *dto.MessagesPager, opt ...options.Option) (
		[]*dto.Message, error)

	// PostMessage 发送频道消息
	PostMessage(ctx context.Context, channelID string, msg *dto.MessageToCreate, opt ...options.Option) (
		*dto.Message, error)

	// PatchMessage 修改频道消息
	PatchMessage(ctx context.Context,
		channelID string, messageID string, msg *dto.MessageToCreate, opt ...options.Option) (*dto.Message, error)

	// RetractMessage 撤回频道消息
	RetractMessage(ctx context.Context, channelID, msgID string, opt ...options.Option) error

	// PostSettingGuide 发送设置引导
	PostSettingGuide(ctx context.Context, channelID string, atUserIDs []string, opt ...options.Option) (
		*dto.Message, error)

	// PostGroupMessage 发送群消息
	PostGroupMessage(ctx context.Context, groupID string, msg dto.APIMessage, opt ...options.Option) (
		*dto.Message, error)

	// PostC2CMessage 发送C2C消息
	PostC2CMessage(ctx context.Context, userID string, msg dto.APIMessage, opt ...options.Option) (
		*dto.Message, error)

	// RetractC2CMessage 撤回C2C消息
	RetractC2CMessage(ctx context.Context, userID, msgID string, opt ...options.Option) error

	// RetractGroupMessage 撤回群消息
	RetractGroupMessage(ctx context.Context, groupID, msgID string, opt ...options.Option) error
}

// GuildAPI guild 相关接口
type GuildAPI interface {
	Guild(ctx context.Context, guildID string) (*dto.Guild, error)
	GuildMember(ctx context.Context, guildID, userID string) (*dto.Member, error)
	GuildMembers(ctx context.Context, guildID string, pager *dto.GuildMembersPager) ([]*dto.Member, error)
	GuildRoleMembers(ctx context.Context, guildID string, roleID string, pager *dto.GuildRoleMembersPager) (
		[]*dto.Member, string, error)
	DeleteGuildMember(ctx context.Context, guildID, userID string, opts ...dto.MemberDeleteOption) error
	// 频道禁言
	GuildMute(ctx context.Context, guildID string, mute *dto.UpdateGuildMute) error
}

// ChannelAPI 频道相关接口
type ChannelAPI interface {
	// Channel 拉取指定子频道信息
	Channel(ctx context.Context, channelID string) (*dto.Channel, error)
	// Channels 拉取子频道列表
	Channels(ctx context.Context, guildID string) ([]*dto.Channel, error)
	// PostChannel 创建子频道
	PostChannel(ctx context.Context, guildID string, value *dto.ChannelValueObject) (*dto.Channel, error)
	// PatchChannel 修改子频道
	PatchChannel(ctx context.Context, channelID string, value *dto.ChannelValueObject) (*dto.Channel, error)
	// DeleteChannel 删除指定子频道
	DeleteChannel(ctx context.Context, channelID string) error
	// CreatePrivateChannel 创建私密子频道
	CreatePrivateChannel(ctx context.Context, guildID string, value *dto.ChannelValueObject, userIDs []string) (
		*dto.Channel, error)
	ListVoiceChannelMembers(ctx context.Context, channelID string) ([]*dto.Member, error)
}

// ChannelPermissionsAPI 子频道权限相关接口
type ChannelPermissionsAPI interface {
	// ChannelPermissions 获取指定子频道的权限
	ChannelPermissions(ctx context.Context, channelID, userID string) (*dto.ChannelPermissions, error)
	// PutChannelPermissions 修改指定子频道的权限
	PutChannelPermissions(ctx context.Context, channelID, userID string, p *dto.UpdateChannelPermissions) error
	// ChannelRolesPermissions  获取指定子频道身份组的权限
	ChannelRolesPermissions(ctx context.Context, channelID, roleID string) (*dto.ChannelRolesPermissions, error)
	// PutChannelRolesPermissions 修改指定子频道身份组的权限
	PutChannelRolesPermissions(ctx context.Context, channelID, roleID string, p *dto.UpdateChannelPermissions) error
}

// AudioAPI 音频接口
type AudioAPI interface {
	// PostAudio 执行音频播放，暂停等操作
	PostAudio(ctx context.Context, channelID string, value *dto.AudioControl) (*dto.AudioControl, error)
	// PutMic 机器人上麦
	PutMic(ctx context.Context, channelID string) error
	// DeleteMic 机器人下麦
	DeleteMic(ctx context.Context, channelID string) error
}

// RoleAPI 用户组相关接口
type RoleAPI interface {
	Roles(ctx context.Context, guildID string) (*dto.GuildRoles, error)
	PostRole(ctx context.Context, guildID string, role *dto.Role) (*dto.UpdateResult, error)
	PatchRole(ctx context.Context, guildID string, roleID dto.RoleID, role *dto.Role) (*dto.UpdateResult, error)
	DeleteRole(ctx context.Context, guildID string, roleID dto.RoleID) error
}

// MemberAPI 成员相关接口，添加成员到用户组等
type MemberAPI interface {
	MemberAddRole(ctx context.Context, guildID string, roleID dto.RoleID, userID string,
		value *dto.MemberAddRoleBody) error

	MemberDeleteRole(ctx context.Context, guildID string, roleID dto.RoleID, userID string,
		value *dto.MemberAddRoleBody) error

	// MemberMute 频道指定单个成员禁言
	MemberMute(ctx context.Context, guildID, userID string, mute *dto.UpdateGuildMute) error

	// MultiMemberMute 频道指定批量成员禁言
	MultiMemberMute(ctx context.Context, guildID string, mute *dto.UpdateGuildMute) (
		*dto.UpdateGuildMuteResponse, error)
}

// DirectMessageAPI 信息相关接口
type DirectMessageAPI interface {
	// CreateDirectMessage 创建私信频道
	CreateDirectMessage(ctx context.Context, dm *dto.DirectMessageToCreate,
		opt ...options.Option) (*dto.DirectMessage, error)

	// PostDirectMessage 在私信频道内发消息
	PostDirectMessage(ctx context.Context, dm *dto.DirectMessage,
		msg *dto.MessageToCreate, opt ...options.Option) (*dto.Message, error)

	// RetractDMMessage 撤回私信频道消息
	RetractDMMessage(ctx context.Context, guildID, msgID string, opt ...options.Option) error

	// PostDMSettingGuide 发送私信设置引导, jumpGuildID为设置引导要跳转的频道ID
	PostDMSettingGuide(ctx context.Context, dm *dto.DirectMessage, jumpGuildID string,
		opt ...options.Option) (*dto.Message, error)
}

// AnnouncesAPI 公告相关接口
type AnnouncesAPI interface {
	// CreateChannelAnnounces 创建子频道公告
	CreateChannelAnnounces(ctx context.Context, channelID string, announce *dto.ChannelAnnouncesToCreate) (
		*dto.Announces, error)

	// DeleteChannelAnnounces 删除子频道公告,会校验 messageID 是否匹配
	DeleteChannelAnnounces(ctx context.Context, channelID, messageID string) error

	// CleanChannelAnnounces 删除子频道公告,不校验 messageID
	CleanChannelAnnounces(ctx context.Context, channelID string) error

	// CreateGuildAnnounces 创建频道全局公告
	CreateGuildAnnounces(ctx context.Context, guildID string, announce *dto.GuildAnnouncesToCreate) (
		*dto.Announces, error)

	// DeleteGuildAnnounces 删除频道全局公告
	DeleteGuildAnnounces(ctx context.Context, guildID, messageID string) error

	// CleanGuildAnnounces 删除频道全局公告,不校验 messageID
	CleanGuildAnnounces(ctx context.Context, guildID string) error
}

// ScheduleAPI 日程相关接口
type ScheduleAPI interface {
	// ListSchedules 查询某个子频道下，since开始的当天的日程列表。若since为0，默认返回当天的日程列表
	ListSchedules(ctx context.Context, channelID string, since uint64) ([]*dto.Schedule, error)

	// GetSchedule 获取单个日程信息
	GetSchedule(ctx context.Context, channelID, scheduleID string) (*dto.Schedule, error)

	// CreateSchedule 创建日程
	CreateSchedule(ctx context.Context, channelID string, schedule *dto.Schedule) (*dto.Schedule, error)

	// ModifySchedule 修改日程
	ModifySchedule(ctx context.Context, channelID, scheduleID string, schedule *dto.Schedule) (*dto.Schedule, error)

	// DeleteSchedule 删除日程
	DeleteSchedule(ctx context.Context, channelID, scheduleID string) error
}

// APIPermissionsAPI api 权限相关接口
type APIPermissionsAPI interface {
	// GetAPIPermissions 获取频道可用权限列表
	GetAPIPermissions(ctx context.Context, guildID string) (*dto.APIPermissions, error)

	// RequireAPIPermissions 创建频道 API 接口权限授权链接
	RequireAPIPermissions(ctx context.Context, guildID string, demand *dto.APIPermissionDemandToCreate) (
		*dto.APIPermissionDemand, error)
}

// PinsAPI 精华消息接口
type PinsAPI interface {
	// AddPins 添加精华消息
	AddPins(ctx context.Context, channelID string, messageID string) (*dto.PinsMessage, error)
	// DeletePins 删除精华消息
	DeletePins(ctx context.Context, channelID, messageID string) error
	// CleanPins 清除全部精华消息
	CleanPins(ctx context.Context, channelID string) error
	// GetPins 获取精华消息
	GetPins(ctx context.Context, channelID string) (*dto.PinsMessage, error)
}

// MessageReactionAPI 消息表情表态接口
type MessageReactionAPI interface {
	// CreateMessageReaction 对消息发表表情表态
	CreateMessageReaction(ctx context.Context, channelID, messageID string, emoji dto.Emoji) error
	// DeleteOwnMessageReaction 删除自己的消息表情表态
	DeleteOwnMessageReaction(ctx context.Context, channelID, messageID string, emoji dto.Emoji) error
	// GetMessageReactionUsers 获取消息表情表态用户列表
	GetMessageReactionUsers(ctx context.Context, channelID, messageID string, emoji dto.Emoji,
		pager *dto.MessageReactionPager) (*dto.MessageReactionUsers, error)
}

// InteractionAPI 互动接口
type InteractionAPI interface {
	// PutInteraction 更新互动信息
	PutInteraction(ctx context.Context, interactionID string, body string) error
}

// WebhookAPI Deprecated: webhook无需维护session。http 事件网关相关接口
type WebhookAPI interface {
	CreateSession(ctx context.Context, identity dto.HTTPIdentity) (*dto.HTTPReady, error)
	CheckSessions(ctx context.Context) ([]*dto.HTTPSession, error)
	SessionList(ctx context.Context) ([]*dto.HTTPSession, error)
	RemoveSession(ctx context.Context, sessionID string) error
}

// MessageSettingAPI 频道消息设置接口
type MessageSettingAPI interface {
	GetMessageSetting(ctx context.Context, guildID string) (*dto.MessageSetting, error)
}
