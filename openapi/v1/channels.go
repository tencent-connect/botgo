package v1

import (
	"context"
	"encoding/json"

	"github.com/tencent-connect/botgo/dto"
)

// Channel 拉取指定子频道信息
func (o *openAPI) Channel(ctx context.Context, channelID string) (*dto.Channel, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Channel{}).
		SetPathParam("channel_id", channelID).
		Get(o.getURL(channelURI))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Channel), nil
}

// Channels 拉取子频道列表
func (o *openAPI) Channels(ctx context.Context, guildID string) ([]*dto.Channel, error) {
	resp, err := o.request(ctx).
		SetPathParam("guild_id", guildID).
		Get(o.getURL(channelsURI))
	if err != nil {
		return nil, err
	}

	channels := make([]*dto.Channel, 0)
	if err := json.Unmarshal(resp.Body(), &channels); err != nil {
		return nil, err
	}

	return channels, nil
}

// PostChannel 创建子频道
func (o *openAPI) PostChannel(ctx context.Context,
	guildID string, value *dto.ChannelValueObject) (*dto.Channel, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Channel{}).
		SetPathParam("guild_id", guildID).
		SetBody(value).
		Post(o.getURL(channelsURI))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Channel), nil
}

// PatchChannel 修改子频道
func (o *openAPI) PatchChannel(ctx context.Context,
	channelID string, value *dto.ChannelValueObject) (*dto.Channel, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Channel{}).
		SetPathParam("channel_id", channelID).
		SetBody(value).
		Patch(o.getURL(channelURI))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Channel), nil
}

// DeleteChannel 删除指定子频道
func (o *openAPI) DeleteChannel(ctx context.Context, channelID string) error {
	_, err := o.request(ctx).
		SetResult(dto.Channel{}).
		SetPathParam("channel_id", channelID).
		Delete(o.getURL(channelURI))
	return err
}

// CreatePrivateChannel 创建私密子频道，底层是用的是 PostChannel 能力
// ChannelValueObject 中的 PrivateType 不需要填充，本方法会自动填充
func (o *openAPI) CreatePrivateChannel(ctx context.Context, guildID string, value *dto.ChannelValueObject,
	userIds []string) (*dto.Channel, error) {
	value.PrivateType = dto.ChannelPrivateTypeAdminAndMember
	if len(userIds) != 0 {
		value.PrivateUserIDs = userIds
		value.PrivateType = dto.ChannelPrivateTypeOnlyAdmin
	}
	return o.PostChannel(ctx, guildID, value)
}

// ListVoiceChannelMembers 查询语音子频道成员列表
func (o *openAPI) ListVoiceChannelMembers(ctx context.Context, channelID string) ([]*dto.Member, error) {
	resp, err := o.request(ctx).
		SetPathParam("channel_id", channelID).
		Get(o.getURL(voiceChannelMembersURI))
	if err != nil {
		return nil, err
	}
	members := make([]*dto.Member, 0)
	if err := json.Unmarshal(resp.Body(), &members); err != nil {
		return nil, err
	}
	return members, nil
}
