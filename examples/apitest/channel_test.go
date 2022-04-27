package apitest

import (
	"testing"

	"github.com/tencent-connect/botgo/dto"
)

func TestChannel(t *testing.T) {
	t.Run(
		"channel list", func(t *testing.T) {
			list, err := api.Channels(ctx, testGuildID)
			if err != nil {
				t.Error(err)
			}
			for _, channel := range list {
				t.Logf("%+v", channel)
			}
			t.Logf(api.TraceID())
		},
	)
	t.Run(
		"create and modify and delete", func(t *testing.T) {
			testGuildID = "3326534247441079828"
			channel, err := api.CreatePrivateChannel(
				ctx, testGuildID, &dto.ChannelValueObject{
					Name:     "机器人创建的频道",
					Type:     dto.ChannelTypeText,
					Position: 0,
					ParentID: "0", // 父ID，正常应该找到一个分组ID，如果传0，就不归属在任何一个分组中
				}, []string{testMemberID},
			)
			if err != nil {
				t.Error(err)
			}

			t.Log(channel)
			channelNew, err := api.PatchChannel(
				ctx, channel.ID, &dto.ChannelValueObject{
					Name: "机器人修改的频道-修改",
				},
			)
			if err != nil {
				t.Error(err)
			}
			t.Log(channelNew)
			if channelNew.Name == channel.Name {
				t.Error("channel name not modified")
			}
			err = api.DeleteChannel(ctx, channelNew.ID)
			if err != nil {
				t.Error(err)
			}
		},
	)
	t.Run(
		"get voice channel member list test", func(t *testing.T) {
			testChannelID := "1572139"
			members, err := api.ListVoiceChannelMembers(ctx, testChannelID)
			if err != nil {
				t.Error(err)
			}
			for _, member := range members {
				t.Logf("member%v", member)
			}

		},
	)
}
