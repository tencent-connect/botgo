package apitest

import (
	"testing"
	"time"

	"github.com/tencent-connect/botgo/dto"
)

func TestAnnounces(t *testing.T) {

	t.Run(
		"create channel announce", func(t *testing.T) {
			messageInfo, err := api.PostMessage(
				ctx, testChannelID, &dto.MessageToCreate{
					Content: "子频道公共创建",
				},
			)
			if err != nil {
				t.Error(err)
			}
			announces, err := api.CreateChannelAnnounces(
				ctx, testChannelID, &dto.ChannelAnnouncesToCreate{
					MessageID: messageInfo.ID,
				},
			)
			if err != nil {
				t.Error(err)
			}
			t.Logf("announces:%+v", announces)
		},
	)
	t.Run(
		"delete channel announce", func(t *testing.T) {
			time.Sleep(3 * time.Second)
			if err := api.DeleteChannelAnnounces(ctx, testChannelID, testMessageID); err != nil {
				t.Error(err)
			}

		},
	)
	t.Run(
		"clean channel announce no check message id", func(t *testing.T) {
			time.Sleep(3 * time.Second)
			err := api.CleanChannelAnnounces(ctx, testChannelID)
			if err != nil {
				t.Error(err)
			}
		},
	)
	t.Run(
		"create guild announce", func(t *testing.T) {
			time.Sleep(3 * time.Second)
			announces, err := api.CreateGuildAnnounces(
				ctx, testGuildID, &dto.GuildAnnouncesToCreate{
					MessageID: testMessageID,
					ChannelID: testChannelID,
				},
			)
			if err != nil {
				t.Error(err)
			}
			t.Logf("announces:%+v", announces)
		},
	)
	t.Run(
		"create recommend channel guild announce", func(t *testing.T) {
			time.Sleep(3 * time.Second)
			announces, err := api.CreateGuildAnnounces(
				ctx, testGuildID, &dto.GuildAnnouncesToCreate{
					AnnouncesType: 0,
					RecommendChannels: []dto.RecommendChannel{
						{
							ChannelID: "1146349",
							Introduce: "子频道 1146349  欢迎语",
						},
						{
							ChannelID: "1703191",
							Introduce: "子频道 1703191  欢迎语",
						},
						{
							ChannelID: "2651556",
							Introduce: "子频道 2651556  欢迎语",
						},
					},
				},
			)
			if err != nil {
				t.Error(err)
			}
			t.Logf("announces:%+v", announces)
		},
	)
	t.Run(
		"delete guild announce", func(t *testing.T) {
			time.Sleep(3 * time.Second)
			if err := api.DeleteGuildAnnounces(ctx, testGuildID, testMessageID); err != nil {
				t.Error(err)
			}
		},
	)
	t.Run(
		"clean guild announce no check message id", func(t *testing.T) {
			time.Sleep(3 * time.Second)
			err := api.CleanGuildAnnounces(ctx, testGuildID)
			if err != nil {
				t.Error(err)
			}
		},
	)
}
