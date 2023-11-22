package v1

import (
	"context"

	"github.com/tencent-connect/botgo/dto"
)

// PostAudio AudioAPI 接口实现
func (o openAPI) PostAudio(ctx context.Context, channelID string, value *dto.AudioControl) (*dto.AudioControl, error) {
	// 目前服务端成功不回包
	_, err := o.request(ctx).
		SetResult(dto.Channel{}).
		SetPathParam("channel_id", channelID).
		SetBody(value).
		Post(o.getURL(audioControlURI))
	if err != nil {
		return nil, err
	}

	return value, nil
}
