package apitest

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/tencent-connect/botgo/dto"
)

const (
	manageChannelPermission     = uint64(1) << 1
	defaultRoleTypeChannelAdmin = "5"
	patchRoleModifyName         = "test role modify"
)

// Test_role 身份组相关接口用例
func Test_role(t *testing.T) {
	var roleID dto.RoleID
	var err error

	t.Run(
		"拉取身份组列表", func(t *testing.T) {
			roles, err := api.Roles(ctx, testGuildID)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%+v", roles)
			for _, role := range roles.Roles {
				t.Logf("%+v", role)
			}
		},
	)
	t.Run(
		"创建身份组", func(t *testing.T) {
			postRoleResult, err := api.PostRole(
				ctx, testGuildID, &dto.Role{
					Name:  "test role",
					Color: 4278245297,
					Hoist: 0,
				},
			)
			if err != nil {
				t.Error(err)
			}
			roleID = postRoleResult.RoleID
			t.Logf("postRoleResult: %+v", postRoleResult)
		},
	)
	t.Run(
		"修改身份组", func(t *testing.T) {
			patchRoleResult, err := api.PatchRole(
				ctx, testGuildID, roleID, &dto.Role{
					Name: patchRoleModifyName,
				},
			)
			if err != nil {
				t.Error(err)
			}
			reflect.DeepEqual(patchRoleModifyName, patchRoleResult.Role.Name)
			t.Logf("patchRoleResult: %+v", patchRoleResult)
		},
	)
	t.Run(
		"删除身份组", func(t *testing.T) {
			err = api.DeleteRole(ctx, testGuildID, roleID)
			if err != nil {
				t.Error(err)
			}
			t.Logf("role id : %v, is deleted", roleID)
		},
	)
}

func Test_roleWithMember(t *testing.T) {
	var roleID dto.RoleID

	t.Run(
		"添加人到身份组", func(t *testing.T) {
			members, err := api.GuildMembers(
				ctx, testGuildID, &dto.GuildMembersPager{
					After: "0",
					Limit: "1",
				},
			)
			if err != nil {
				t.Error(err)
			}
			userID := members[0].User.ID
			err = api.MemberAddRole(ctx, testGuildID, roleID, userID, nil)
			if err != nil {
				t.Error(err)
			}
			member, err := api.GuildMember(ctx, testGuildID, userID)
			var roleFound bool
			for _, role := range member.Roles {
				if role == string(roleID) {
					roleFound = true
				}
			}
			if !roleFound {
				t.Error("not found role id been add")
			}
		},
	)
	t.Run(
		"添加人到子频道管理员身份组并指定子频道", func(t *testing.T) {
			members, err := api.GuildMembers(
				ctx, testGuildID, &dto.GuildMembersPager{
					After: "0",
					Limit: "1",
				},
			)
			if err != nil {
				t.Error(err)
			}
			userID := members[0].User.ID
			channels, err := api.Channels(ctx, testGuildID)
			if err != nil {
				t.Error(err)
			}
			channelID := channels[len(channels)-1].ID
			t.Logf("testGuildID: %+v, channelID: %+v", testGuildID, channelID)
			err = api.MemberAddRole(
				ctx, testGuildID, defaultRoleTypeChannelAdmin, userID, &dto.MemberAddRoleBody{
					Channel: &dto.Channel{
						ID: channelID,
					},
				},
			)
			if err != nil {
				t.Error(err)
			}
			channelPermissions, err := api.ChannelPermissions(ctx, channelID, userID)
			if err != nil {
				t.Error(err)
			}
			channelPermissionsUint, err := strconv.ParseUint(channelPermissions.Permissions, 10, 64)
			if err != nil {
				t.Error(err)
			}
			t.Logf(
				"channelPermissionsUint: %+v, channelPermissions.Permissions: %+v",
				channelPermissionsUint, channelPermissions.Permissions,
			)
			if channelPermissionsUint&manageChannelPermission != 2 {
				t.Error("not found channel permissions been add")
			}
		},
	)
	t.Run(
		"在子频道管理员身份组删除用户并指定子频道", func(t *testing.T) {
			members, err := api.GuildMembers(
				ctx, testGuildID, &dto.GuildMembersPager{
					After: "0",
					Limit: "1",
				},
			)
			if err != nil {
				t.Error(err)
			}
			userID := members[0].User.ID
			channels, err := api.Channels(ctx, testGuildID)
			if err != nil {
				t.Error(err)
			}
			channelID := channels[len(channels)-1].ID
			t.Logf("testGuildID: %+v, channelID: %+v", testGuildID, channelID)
			err = api.MemberDeleteRole(
				ctx, testGuildID, defaultRoleTypeChannelAdmin, userID, &dto.MemberAddRoleBody{
					Channel: &dto.Channel{
						ID: channelID,
					},
				},
			)
			if err != nil {
				t.Error(err)
			}
			channelPermissions, err := api.ChannelPermissions(ctx, channelID, userID)
			if err != nil {
				t.Error(err)
			}
			channelPermissionsUint, err := strconv.ParseUint(channelPermissions.Permissions, 10, 64)
			if err != nil {
				t.Error(err)
			}
			t.Logf("channelPermissionsUint: %+v", channelPermissionsUint)
			if channelPermissionsUint&manageChannelPermission == 2 {
				t.Error("not found channel permissions been add")
			}
		},
	)
}
