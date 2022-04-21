package v1

import (
	"context"
	"errors"

	"github.com/tencent-connect/botgo/log"

	"github.com/tencent-connect/botgo/dto"
)

// GuildMute 频道禁言
func (o *openAPI) GuildMute(ctx context.Context, guildID string, mute *dto.UpdateGuildMute) error {
	_, err := o.request(ctx).
		SetPathParam("guild_id", guildID).
		SetBody(mute).
		Patch(o.getURL(guildMuteURI))
	if err != nil {
		return err
	}
	return nil
}

// MemberMute 频道指定成员禁言
func (o *openAPI) MemberMute(ctx context.Context, guildID, userID string,
	mute *dto.UpdateGuildMute) error {
	_, err := o.request(ctx).
		SetPathParam("guild_id", guildID).
		SetPathParam("user_id", userID).
		SetBody(mute).
		Patch(o.getURL(guildMembersMuteURI))
	if err != nil {
		return err
	}
	return nil
}

// MultiMemberMute 频道批量成员禁言
func (o *openAPI) MultiMemberMute(ctx context.Context, guildID string,
	mute *dto.UpdateGuildMute) (*dto.UpdateGuildMuteResponse, error) {
	if len(mute.UserIDs) == 0 {
		return nil, errors.New("no user id param")
	}
	rsp, err := o.request(ctx).
		SetPathParam("guild_id", guildID).
		SetBody(mute).
		SetResult(dto.UpdateGuildMuteResponse{}).
		Patch(o.getURL(guildMuteURI))
	if err != nil {
		return nil, err
	}
	log.Infof("MultiMemberMute rsp result: %#v", rsp.Result())
	return rsp.Result().(*dto.UpdateGuildMuteResponse), nil
}
