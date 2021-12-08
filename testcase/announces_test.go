package testcase

import (
	"testing"

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
		t.Logf("announces:%+v",announces)
	})
	t.Run("delete channel announce", func(t *testing.T) {
		err := api.DeleteChannelAnnounces(ctx, testChannelID, testMessageID)
		if err != nil {
			t.Error(err)
		}
	})
}
