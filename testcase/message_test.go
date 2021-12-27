package testcase

import (
	"testing"

	"github.com/tencent-connect/botgo/dto"
)

func TestMessage(t *testing.T) {
	t.Run("message list", func(t *testing.T) {
		// 先拉取3条消息
		messages, err := api.Messages(ctx, testChannelID, &dto.MessagesPager{
			Limit: "3",
		})
		if err != nil {
			t.Error(err)
		}
		index := make(map[int]string)
		for i, message := range messages {
			index[i] = message.ID
			t.Log(message.ID, message.Author.Username, message.Timestamp)
		}

		// 从上面3条的第二条往前拉取
		messages, err = api.Messages(ctx, testChannelID, &dto.MessagesPager{
			Type:  dto.MPTBefore,
			ID:    index[1],
			Limit: "2",
		})
		if err != nil {
			t.Error(err)
		}
		for i, message := range messages {
			if i == 2 && index[2] != message.ID {
				t.Error("before id not match")
			}
			t.Log(message.ID, message.Author.Username, message.Timestamp)
		}

		// 从上面3条的第二条往后拉取
		messages, err = api.Messages(ctx, testChannelID, &dto.MessagesPager{
			Type:  dto.MPTAfter,
			ID:    index[1],
			Limit: "2",
		})
		if err != nil {
			t.Error(err)
		}
		for i, message := range messages {
			if i == 0 && index[0] != message.ID {
				t.Error("after id not match")
			}
			t.Log(message.ID, message.Author.Username, message.Timestamp)
		}
		// 从上面3条的第二条环绕拉取
		messages, err = api.Messages(ctx, testChannelID, &dto.MessagesPager{
			Type:  dto.MPTAround,
			ID:    index[1],
			Limit: "3",
		})
		if err != nil {
			t.Error(err)
		}
		for i, message := range messages {
			if i == 0 && index[0] != message.ID {
				t.Error("around id not match")
			}
			if i == 2 && index[2] != message.ID {
				t.Error("around id not match")
			}
			t.Log(message.ID, message.Author.Username, message.Timestamp)
		}
	})
}

func TestRetractMessage(t *testing.T) {
	msgID := "109b8a401a1231343431313532313831383136323933383420801e28003081b0f30338cd6040c36048f5e4908e0650b1acf8fa05"
	t.Run("消息撤回", func(t *testing.T) {
		err := api.RetractMessage(ctx, "1049883", msgID)
		if err != nil {
			t.Error(err)
		}
		t.Logf("msg id : %v, is deleted", msgID)
	})
}
