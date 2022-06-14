package v1

import (
	"context"
	"encoding/json"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/errs"
)

// MemberAddRole 添加成员角色
func (o *openAPI) MemberAddRole(
	ctx context.Context, guildID string, roleID dto.RoleID, userID string,
	value *dto.MemberAddRoleBody,
) error {
	if value == nil {
		value = new(dto.MemberAddRoleBody)
	}
	_, err := o.request(ctx).
		SetPathParam("guild_id", guildID).
		SetPathParam("role_id", string(roleID)).
		SetPathParam("user_id", userID).
		SetBody(value).
		Put(o.getURL(memberRoleURI))
	return err
}

// MemberDeleteRole 删除成员角色
func (o *openAPI) MemberDeleteRole(
	ctx context.Context, guildID string, roleID dto.RoleID, userID string,
	value *dto.MemberAddRoleBody,
) error {
	if value == nil {
		value = new(dto.MemberAddRoleBody)
	}
	_, err := o.request(ctx).
		SetPathParam("guild_id", guildID).
		SetPathParam("role_id", string(roleID)).
		SetPathParam("user_id", userID).
		SetBody(value).
		Delete(o.getURL(memberRoleURI))
	return err
}

// GuildMember 拉取频道指定成员
func (o *openAPI) GuildMember(ctx context.Context, guildID, userID string) (*dto.Member, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Member{}).
		SetPathParam("guild_id", guildID).
		SetPathParam("user_id", userID).
		Get(o.getURL(guildMemberURI))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Member), nil
}

// GuildMembers 分页拉取频道内成员列表
func (o *openAPI) GuildMembers(
	ctx context.Context,
	guildID string, pager *dto.GuildMembersPager,
) ([]*dto.Member, error) {
	if pager == nil {
		return nil, errs.ErrPagerIsNil
	}
	resp, err := o.request(ctx).
		SetPathParam("guild_id", guildID).
		SetQueryParams(pager.QueryParams()).
		Get(o.getURL(guildMembersURI))
	if err != nil {
		return nil, err
	}

	members := make([]*dto.Member, 0)
	if err := json.Unmarshal(resp.Body(), &members); err != nil {
		return nil, err
	}

	return members, nil
}

// GuildRoleMembers 分页拉取频道内身份组成员列表
func (o *openAPI) GuildRoleMembers(
	ctx context.Context, guildID string, roleID string, pager *dto.GuildRoleMembersPager,
) ([]*dto.Member, string, error) {
	if pager == nil {
		return nil, "", errs.ErrPagerIsNil
	}
	resp, err := o.request(ctx).
		SetPathParam("guild_id", guildID).
		SetPathParam("role_id", roleID).
		SetQueryParams(pager.QueryParams()).
		Get(o.getURL(guildRoleMemberURI))
	if err != nil {
		return nil, "", err
	}

	type res struct {
		Data []*dto.Member `json:"data"`
		Next string        `json:"next"`
	}
	var roleMembersRsp res
	if err := json.Unmarshal(resp.Body(), &roleMembersRsp); err != nil {
		return nil, "", err
	}

	return roleMembersRsp.Data, roleMembersRsp.Next, nil
}

// DeleteGuildMember 将指定成员踢出频道
func (o *openAPI) DeleteGuildMember(ctx context.Context, guildID, userID string, opts ...dto.MemberDeleteOption) error {
	opt := &dto.MemberDeleteOpts{}
	for _, o := range opts {
		o(opt)
	}
	_, err := o.request(ctx).
		SetResult(dto.Member{}).
		SetPathParam("guild_id", guildID).
		SetPathParam("user_id", userID).
		SetBody(opt).
		Delete(o.getURL(guildMemberURI))
	return err
}
