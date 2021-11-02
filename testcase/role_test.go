package testcase

import (
	"testing"

	"github.com/tencent-connect/botgo/dto"
)

// Test_role 用户组相关接口用例
func Test_role(t *testing.T) {
	var roleID dto.RoleID
	var err error

	t.Run("拉取用户组列表", func(t *testing.T) {
		roles, err := api.Roles(ctx, testGuildID)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%+v", roles)
		for _, role := range roles.Roles {
			t.Logf("%+v", role)
		}
	})
	t.Run("创建用户组", func(t *testing.T) {
		roleID, err = api.PostRole(ctx, testGuildID, &dto.Role{
			Name:  "test role",
			Color: 4278245297,
			Hoist: 0,
		})
		if err != nil {
			t.Error(err)
		}
		t.Logf("role id : %v", roleID)
	})
	t.Run("添加人到用户组", func(t *testing.T) {
		members, err := api.GuildMembers(ctx, testGuildID, &dto.GuildMembersPager{
			After: "0",
			Limit: "1",
		})
		if err != nil {
			t.Error(err)
		}
		userID := members[0].User.ID
		err = api.MemberAddRole(ctx, testGuildID, roleID, userID)
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
	})
	t.Run("删除用户组", func(t *testing.T) {
		err = api.DeleteRole(ctx, testGuildID, roleID)
		if err != nil {
			t.Error(err)
		}
		t.Logf("role id : %v, is deleted", roleID)
	})
}
