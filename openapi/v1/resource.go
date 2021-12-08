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
	guildURI        uri = "/guilds/{guild_id}"
	guildMembersURI uri = "/guilds/{guild_id}/members"
	guildMemberURI  uri = "/guilds/{guild_id}/members/{user_id}"

	channelsURI uri = "/guilds/{guild_id}/channels"
	channelURI  uri = "/channels/{channel_id}"

	channelPermissionsURI uri = "/channels/{channel_id}/members/{user_id}/permissions"

	messagesURI uri = "/channels/{channel_id}/messages"
	messageURI  uri = "/channels/{channel_id}/messages/{message_id}"

	userMeURI       uri = "/users/@me"
	userMeGuildsURI uri = "/users/@me/guilds"
	userMeDMURI     uri = "/users/@me/dms"

	gatewayURI    uri = "/gateway" // nolint
	gatewayBotURI uri = "/gateway/bot"

	audioControlURI uri = "/channels/{channel_id}/audio"

	rolesURI uri = "/guilds/{guild_id}/roles"
	roleURI  uri = "/guilds/{guild_id}/roles/{role_id}"

	memberRoleURI uri = "/guilds/{guild_id}/members/{user_id}/roles/{role_id}"

	dmsURI uri = "/dms/{guild_id}/messages"

	announcesURI = "/channels/{channel_id}/announces"
	announceURI  = "/channels/{channel_id}/announces/{message_id}"
)

func getURL(endpoint uri, sandbox bool) string {
	d := domain
	if sandbox {
		d = sandBoxDomain
	}
	return fmt.Sprintf("%s://%s%s", scheme, d, endpoint)
}
