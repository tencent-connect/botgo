package testcase

import (
	"testing"

	"github.com/tencent-connect/botgo/dto"
)

func TestChannel(t *testing.T) {
	t.Run("guild info", func(t *testing.T) {
		guild, err := api.Guild(ctx, testGuildID)
		if err != nil {
			t.Error(err)
		}
		t.Log(guild)
	})
	t.Run("channel list", func(t *testing.T) {
		list, err := api.Channels(ctx, testGuildID)
		if err != nil {
			t.Error(err)
		}
		for _, channel := range list {
			t.Logf("%+v", channel)
		}
		t.Logf(api.TraceID())
	})
	t.Run("create live channel", func(t *testing.T) {
		api.PostChannel(ctx, testGuildID, &dto.ChannelValueObject{
			Name:     "机器人创建2",
			Type:     dto.ChannelTypeLive,
			Position: 0,   // 默认是当前时间戳，如果传递，则要避免和其他频道的 position 重复，否则会报错
			ParentID: "0", // 父ID，正常应该找到一个分组ID，如果传0，就不归属在任何一个分组中
		})
	})
}
