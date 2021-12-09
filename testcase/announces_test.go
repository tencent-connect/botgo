package testcase

import (
	"testing"
	"time"

	"github.com/tencent-connect/botgo/dto"
)

func TestAnnounces(t *testing.T) {
	var messageID string
	t.Run("create channel announce", func(t *testing.T) {
		messageInfo, err := api.PostMessage(ctx, testChannelID, &dto.MessageToCreate{
			Content: "子频道公共创建",
		})
		if err != nil {
			t.Error(err)
		}
		messageID = messageInfo.ID
		announces, err := api.CreateChannelAnnounces(ctx, testChannelID, &dto.ChannelAnnouncesToCreate{
			MessageID: messageID,
		})
		if err != nil {
			t.Error(err)
		}
		t.Logf("announces:%+v", announces)
	})
	t.Run("delete channel announce", func(t *testing.T) {
		time.Sleep(3 * time.Second)
		err := api.DeleteChannelAnnounces(ctx, testChannelID, messageID)
		if err != nil {
			t.Error(err)
		}
	})
}
