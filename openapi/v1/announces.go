package v1

import (
	"context"

	"github.com/tencent-connect/botgo/dto"
)

// CreateChannelAnnounces 创子频道公告
func (o *openAPI) CreateChannelAnnounces(ctx context.Context, channelID string,
	announce *dto.ChannelAnnouncesToCreate) (*dto.Announces, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Announces{}).
		SetPathParam("channel_id", channelID).
		SetBody(announce).
		Post(getURL(announcesURI, o.sandbox))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.Announces), nil
}

// DeleteChannelAnnounces 删除子频道公告
func (o *openAPI) DeleteChannelAnnounces(ctx context.Context, channelID, messageID string) error {
	_, err := o.request(ctx).
		SetResult(dto.Announces{}).
		SetPathParam("channel_id", channelID).
		SetPathParam("message_id", messageID).
		Delete(getURL(announceURI, o.sandbox))
	return err
}
