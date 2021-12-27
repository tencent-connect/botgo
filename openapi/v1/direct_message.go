package v1

import (
	"context"

	"github.com/tencent-connect/botgo/dto"
)

// CreateDirectMessage 创建私信频道
func (o *openAPI) CreateDirectMessage(ctx context.Context, dm *dto.DirectMessageToCreate) (*dto.DirectMessage, error) {
	resp, err := o.request(ctx).
		SetResult(dto.DirectMessage{}).
		SetBody(dm).
		Post(getURL(userMeDMURI, o.sandbox))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.DirectMessage), nil
}

// PostDirectMessage 在私信频道内发消息
func (o *openAPI) PostDirectMessage(ctx context.Context,
	dm *dto.DirectMessage, msg *dto.MessageToCreate) (*dto.Message, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("guild_id", dm.GuildID).
		SetBody(msg).
		Post(getURL(dmsURI, o.sandbox))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.Message), nil
}

// RetractDMMessage 撤回私信消息
func (o *openAPI) RetractDMMessage(ctx context.Context, guildID, msgID string) error {
	_, err := o.request(ctx).
		SetPathParam("guild_id", guildID).
		SetPathParam("message_id", string(msgID)).
		Delete(getURL(dmsMessageURI, o.sandbox))
	return err
}
