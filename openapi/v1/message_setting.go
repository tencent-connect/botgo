package v1

import (
	"context"

	"github.com/tencent-connect/botgo/dto"
)

// GetMessageSetting 获取频道消息频率设置信息
func (o *openAPI) GetMessageSetting(ctx context.Context, guildID string) (*dto.MessageSetting, error) {
	resp, err := o.request(ctx).
		SetResult(dto.MessageSetting{}).
		SetPathParam("guild_id", guildID).
		Get(o.getURL(messageSettingURI))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.MessageSetting), nil
}
