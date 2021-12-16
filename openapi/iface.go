package openapi

import (
	"context"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/token"
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
}

// Base 基础能力接口
type Base interface {
	Version() APIVersion
	New(token *token.Token, inSandbox bool) OpenAPI
	WithTimeout(duration time.Duration) OpenAPI
	// WithBody 设置 body，如果 openapi 提供设置 body 的功能，则需要自行识别 body 类型
	WithBody(body interface{}) OpenAPI
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
	Message(ctx context.Context, channelID string, messageID string) (*dto.Message, error)
	Messages(ctx context.Context, channelID string, pager *dto.MessagesPager) ([]*dto.Message, error)
	PostMessage(ctx context.Context, channelID string, msg *dto.MessageToCreate) (*dto.Message, error)
}

// GuildAPI guild 相关接口
type GuildAPI interface {
	Guild(ctx context.Context, guildID string) (*dto.Guild, error)
	GuildMember(ctx context.Context, guildID, userID string) (*dto.Member, error)
	GuildMembers(ctx context.Context, guildID string, pager *dto.GuildMembersPager) ([]*dto.Member, error)
	DeleteGuildMember(ctx context.Context, guildID, userID string) error
}

// ChannelAPI 频道相关接口
type ChannelAPI interface {
	Channel(ctx context.Context, channelID string) (*dto.Channel, error)
	Channels(ctx context.Context, guildID string) ([]*dto.Channel, error)
	PostChannel(ctx context.Context, guildID string, value *dto.ChannelValueObject) (*dto.Channel, error)
	PatchChannel(ctx context.Context, channelID string, value *dto.ChannelValueObject) (*dto.Channel, error)
	DeleteChannel(ctx context.Context, channelID string) error
}

// ChannelPermissionsAPI 子频道权限相关接口
type ChannelPermissionsAPI interface {
	// ChannelPermissions 获取指定子频道的权限
	ChannelPermissions(ctx context.Context, channelID, userID string) (*dto.ChannelPermissions, error)
	// PutChannelPermissions 修改指定子频道的权限
	PutChannelPermissions(ctx context.Context, channelID, userID string, p *dto.UpdateChannelPermissions) error
}

// AudioAPI 音频接口
type AudioAPI interface {
	// PostAudio 执行音频播放，暂停等操作
	PostAudio(ctx context.Context, channelID string, value *dto.AudioControl) (*dto.AudioControl, error)
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
	MemberAddRole(ctx context.Context,
		guildID string, roleID dto.RoleID, userID string, value *dto.MemberAddRoleBody) error
	MemberDeleteRole(ctx context.Context,
		guildID string, roleID dto.RoleID, userID string, value *dto.MemberAddRoleBody) error
}

// DirectMessageAPI 信息相关接口
type DirectMessageAPI interface {
	// CreateDirectMessage 创建私信频道
	CreateDirectMessage(ctx context.Context, dm *dto.DirectMessageToCreate) (*dto.DirectMessage, error)
	// PostDirectMessage 在私信频道内发消息
	PostDirectMessage(ctx context.Context, dm *dto.DirectMessage, msg *dto.MessageToCreate) (*dto.Message, error)
}

// AnnouncesAPI 公告相关接口
type AnnouncesAPI interface {
	// CreateChannelAnnounces 创建子频道公告
	CreateChannelAnnounces(ctx context.Context,
		channelID string, announce *dto.ChannelAnnouncesToCreate) (*dto.Announces, error)
	// DeleteChannelAnnounces 删除子频道公告,会校验 messageID 是否匹配
	DeleteChannelAnnounces(ctx context.Context, channelID, messageID string) error
	// CleanChannelAnnounces 删除子频道公告,不校验 messageID
	CleanChannelAnnounces(ctx context.Context, channelID string) error
	// CreateGuildAnnounces 创建频道全局公告
	CreateGuildAnnounces(ctx context.Context, guildID string,
		announce *dto.GuildAnnouncesToCreate) (*dto.Announces, error)
	// DeleteGuildAnnounces 删除频道全局公告
	DeleteGuildAnnounces(ctx context.Context, guildID, messageID string) error
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
