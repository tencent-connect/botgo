package v1

import (
	"fmt"
)

const domain = "api.sgroup.qq.com"
const sandBoxDomain = "sandbox.api.sgroup.qq.com"

const scheme = "https"

type uri string

// 目前提供的接口的 uri
const (
	guildURI            uri = "/guilds/{guild_id}"
	guildMembersURI     uri = "/guilds/{guild_id}/members"
	guildMemberURI      uri = "/guilds/{guild_id}/members/{user_id}"
	guildRoleMemberURI  uri = "/guilds/{guild_id}/roles/{role_id}/members"
	guildMuteURI        uri = "/guilds/{guild_id}/mute"                   // 频道禁言
	guildMembersMuteURI uri = "/guilds/{guild_id}/members/{user_id}/mute" // 频道指定成员禁言

	channelsURI uri = "/guilds/{guild_id}/channels"
	channelURI  uri = "/channels/{channel_id}"

	channelPermissionsURI      uri = "/channels/{channel_id}/members/{user_id}/permissions"
	channelRolesPermissionsURI uri = "/channels/{channel_id}/roles/{role_id}/permissions"

	messagesURI uri = "/channels/{channel_id}/messages"
	messageURI  uri = "/channels/{channel_id}/messages/{message_id}"

	userMeURI       uri = "/users/@me"
	userMeGuildsURI uri = "/users/@me/guilds"
	userMeDMURI     uri = "/users/@me/dms"

	gatewayURI    uri = "/gateway" // nolint
	gatewayBotURI uri = "/gateway/bot"

	audioControlURI uri = "/channels/{channel_id}/audio"
	micURI          uri = "/channels/{channel_id}/mic"

	rolesURI uri = "/guilds/{guild_id}/roles"
	roleURI  uri = "/guilds/{guild_id}/roles/{role_id}"

	memberRoleURI uri = "/guilds/{guild_id}/members/{user_id}/roles/{role_id}"

	dmsURI        uri = "/dms/{guild_id}/messages"
	dmsMessageURI uri = "/dms/{guild_id}/messages/{message_id}"

	channelAnnouncesURI = "/channels/{channel_id}/announces"
	channelAnnounceURI  = "/channels/{channel_id}/announces/{message_id}"

	guildAnnouncesURI = "/guilds/{guild_id}/announces"
	guildAnnounceURI  = "/guilds/{guild_id}/announces/{message_id}"

	schedulesURI uri = "/channels/{channel_id}/schedules"
	scheduleURI  uri = "/channels/{channel_id}/schedules/{schedule_id}"

	apiPermissionURI       uri = "/guilds/{guild_id}/api_permission"
	apiPermissionDemandURI uri = "/guilds/{guild_id}/api_permission/demand"

	pinsURI = "/channels/{channel_id}/pins"
	pinURI  = "/channels/{channel_id}/pins/{message_id}"

	messageReactionURI uri = "/channels/{channel_id}/messages/{message_id}/reactions/{emoji_type}/{emoji_id}"

	interactionsURI = "/interactions/{interaction_id}"

	httpSessionsURI uri = "/gateway/webhook/sessions"
	httpSessionURI  uri = "/gateway/webhook/sessions/{session_id}"

	messageSettingURI uri = "/guilds/{guild_id}/message/setting"

	voiceChannelMembersURI uri = "/channels/{channel_id}/voice/members"

	settingGuideURI   uri = "/channels/{channel_id}/settingguide"
	dmSettingGuideURI uri = "/dms/{guild_id}/settingguide"
)

// getURL 获取接口地址，会处理沙箱环境判断
func (o *openAPI) getURL(endpoint uri) string {
	d := domain
	if o.sandbox {
		d = sandBoxDomain
	}
	return fmt.Sprintf("%s://%s%s", scheme, d, endpoint)
}
