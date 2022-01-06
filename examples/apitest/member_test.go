package apitest

import (
	"testing"

	"github.com/tencent-connect/botgo/dto"
)

func Test_Member(t *testing.T) {
	var userId string
	t.Run(
		"list member", func(t *testing.T) {
			members, err := api.GuildMembers(
				ctx, testGuildID, &dto.GuildMembersPager{
					After: "0",
					Limit: "10",
				},
			)
			for _, member := range members {
				t.Logf("user: %+v", member.User.Username)
				t.Logf("roles: %+v", member.Roles)
				userId = member.User.ID
			}
			if err != nil {
				t.Error(err)
			}
		},
	)
	t.Run(
		"get member", func(t *testing.T) {
			member, err := api.GuildMember(ctx, testGuildID, userId)
			if err != nil {
				t.Error(err)
			}
			t.Logf("member: %+v", member)
			t.Logf("user: %+v", member.User)
		},
	)
}
