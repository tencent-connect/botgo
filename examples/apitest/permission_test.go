package apitest

import (
	"testing"

	"github.com/tencent-connect/botgo/dto"
)

func Test_ChannelRolesPermissions(t *testing.T) {
	t.Run(
		"update roles Permissions", func(t *testing.T) {
			updatePermissions := &dto.UpdateChannelPermissions{
				Add: "5",
			}
			err := api.PutChannelRolesPermissions(ctx, testChannelID, testRolesID, updatePermissions)
			if err != nil {
				t.Error(err)
			}
		},
	)
	t.Run(
		"get roles Permissions", func(t *testing.T) {
			permissions, err := api.ChannelRolesPermissions(ctx, testChannelID, testRolesID)
			if err != nil {
				t.Error(err)
			}
			t.Logf("permissions: %+v", permissions)
		},
	)
}

func Test_ChannelMembersPermissions(t *testing.T) {
	t.Run(
		"update members Permissions", func(t *testing.T) {
			updatePermissions := &dto.UpdateChannelPermissions{
				Add: "5",
			}
			err := api.PutChannelPermissions(ctx, testChannelID, testMemberID, updatePermissions)
			if err != nil {
				t.Error(err)
			}
		},
	)
	t.Run(
		"get members Permissions", func(t *testing.T) {
			permissions, err := api.ChannelPermissions(ctx, testChannelID, testMemberID)
			if err != nil {
				t.Error(err)
			}
			t.Logf("permissions: %+v", permissions)
		},
	)
}
