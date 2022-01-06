package apitest

import (
	"testing"

	"github.com/tencent-connect/botgo/dto"
)

func TestGuild(t *testing.T) {
	t.Run(
		"guild info", func(t *testing.T) {
			guild, err := api.Guild(ctx, testGuildID)
			if err != nil {
				t.Error(err)
			}
			t.Log(guild)
		},
	)
	t.Run(
		"my join guilds", func(t *testing.T) {
			guilds, err := api.MeGuilds(
				ctx, &dto.GuildPager{
					Limit: "100",
				},
			)
			if err != nil {
				t.Error(err)
			}
			for _, guild := range guilds {
				t.Log(guild)
			}
		},
	)
}
