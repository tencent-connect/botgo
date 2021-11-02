package v1

import (
	"context"

	"github.com/tencent-connect/botgo/dto"
)

func (o *openAPI) MemberAddRole(ctx context.Context, guildID string, roleID dto.RoleID, userID string) error {
	_, err := o.request(ctx).
		SetPathParam("guild_id", guildID).
		SetPathParam("role_id", string(roleID)).
		SetPathParam("user_id", userID).
		Put(getURL(memberRoleURI, o.sandbox))
	return err
}

func (o *openAPI) MemberDeleteRole(ctx context.Context, guildID string, roleID dto.RoleID, userID string) error {
	_, err := o.request(ctx).
		SetPathParam("guild_id", guildID).
		SetPathParam("role_id", string(roleID)).
		SetPathParam("user_id", userID).
		Delete(getURL(memberRoleURI, o.sandbox))
	return err
}
