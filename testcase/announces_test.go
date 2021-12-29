package testcase

import (
	"testing"
	"time"

	"github.com/tencent-connect/botgo/dto"
)

func TestAnnounces(t *testing.T) {

	t.Run("create channel announce", func(t *testing.T) {
		messageInfo, err := api.PostMessage(ctx, testChannelID, &dto.MessageToCreate{
			Content: "子频道公共创建",
		})
		if err != nil {
			t.Error(err)
		}
		announces, err := api.CreateChannelAnnounces(ctx, testChannelID, &dto.ChannelAnnouncesToCreate{
			MessageID: messageInfo.ID,
		})
		if err != nil {
			t.Error(err)
		}
		t.Logf("announces:%+v", announces)
	})
	t.Run("delete channel announce", func(t *testing.T) {
		time.Sleep(3 * time.Second)
		if err := api.DeleteChannelAnnounces(ctx, testChannelID, testMessageID); err != nil {
			t.Error(err)
		}

	})
	t.Run("clean channel announce no check messageID", func(t *testing.T) {
		time.Sleep(3 * time.Second)
		err := api.CleanChannelAnnounces(ctx, testChannelID)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("create guild announce", func(t *testing.T) {
		time.Sleep(3 * time.Second)
		announces, err := api.CreateGuildAnnounces(ctx, testGuildID, &dto.GuildAnnouncesToCreate{
			MessageID: testMessageID,
			ChannelID: testChannelID,
		})
		if err != nil {
			t.Error(err)
		}
		t.Logf("announces:%+v", announces)
	})
	t.Run("delete guild announce", func(t *testing.T) {
		time.Sleep(3 * time.Second)
		if err := api.DeleteGuildAnnounces(ctx, testGuildID, testMessageID); err != nil {
			t.Error(err)
		}
	})
	t.Run("clean guild announce no check messageID", func(t *testing.T) {
		time.Sleep(3 * time.Second)
		err := api.CleanGuildAnnounces(ctx, testGuildID)
		if err != nil {
			t.Error(err)
		}
	})
}
