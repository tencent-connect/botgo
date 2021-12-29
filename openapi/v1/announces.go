package v1

import (
	"context"

	"github.com/tencent-connect/botgo/dto"
)

// CreateChannelAnnounces 创建子频道公告
func (o *openAPI) CreateChannelAnnounces(ctx context.Context, channelID string,
	announce *dto.ChannelAnnouncesToCreate) (*dto.Announces, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Announces{}).
		SetPathParam("channel_id", channelID).
		SetBody(announce).
		Post(getURL(channelAnnouncesURI, o.sandbox))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.Announces), nil
}

// DeleteChannelAnnounces 删除子频道公告,会校验 messageID
func (o *openAPI) DeleteChannelAnnounces(ctx context.Context, channelID, messageID string) error {
	_, err := o.request(ctx).
		SetResult(dto.Announces{}).
		SetPathParam("channel_id", channelID).
		SetPathParam("message_id", messageID).
		Delete(getURL(channelAnnounceURI, o.sandbox))
	return err
}

// CleanChannelAnnounces 删除子频道公告,不校验 messageID
func (o *openAPI) CleanChannelAnnounces(ctx context.Context, channelID string) error {
	_, err := o.request(ctx).
		SetResult(dto.Announces{}).
		SetPathParam("channel_id", channelID).
		SetPathParam("message_id", "all").
		Delete(getURL(channelAnnounceURI, o.sandbox))
	return err
}

// CreateGuildAnnounces 创建频道全局公告
func (o *openAPI) CreateGuildAnnounces(ctx context.Context, guildID string,
	announce *dto.GuildAnnouncesToCreate) (*dto.Announces, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Announces{}).
		SetPathParam("guild_id", guildID).
		SetBody(announce).
		Post(getURL(guildAnnouncesURI, o.sandbox))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.Announces), nil
}

// DeleteGuildAnnounces 删除频道全局公告,会校验 messageID
func (o *openAPI) DeleteGuildAnnounces(ctx context.Context, guildID, messageID string) error {
	_, err := o.request(ctx).
		SetResult(dto.Announces{}).
		SetPathParam("guild_id", guildID).
		SetPathParam("message_id", messageID).
		Delete(getURL(guildAnnounceURI, o.sandbox))
	return err
}

// CleanGuildAnnounces 删除道全局公告,不校验 messageID
func (o *openAPI) CleanGuildAnnounces(ctx context.Context, guildID string) error {
	_, err := o.request(ctx).
		SetResult(dto.Announces{}).
		SetPathParam("guild_id", guildID).
		SetPathParam("message_id", "all").
		Delete(getURL(guildAnnounceURI, o.sandbox))
	return err
}
