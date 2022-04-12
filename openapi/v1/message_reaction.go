package v1

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/errs"
)

// CreateMessageReaction 对消息发表表情表态
func (o *openAPI) CreateMessageReaction(ctx context.Context,
	channelID, messageID string, emoji dto.Emoji) error {
	_, err := o.request(ctx).
		SetPathParam("channel_id", channelID).
		SetPathParam("message_id", messageID).
		SetPathParam("emoji_type", strconv.FormatUint(uint64(emoji.Type), 10)).
		SetPathParam("emoji_id", emoji.ID).
		Put(o.getURL(messageReactionURI))
	if err != nil {
		return err
	}
	return nil
}

// DeleteOwnMessageReaction 删除自己的消息表情表态
func (o *openAPI) DeleteOwnMessageReaction(ctx context.Context,
	channelID, messageID string, emoji dto.Emoji) error {
	_, err := o.request(ctx).
		SetPathParam("channel_id", channelID).
		SetPathParam("message_id", messageID).
		SetPathParam("emoji_type", strconv.FormatUint(uint64(emoji.Type), 10)).
		SetPathParam("emoji_id", emoji.ID).
		Delete(o.getURL(messageReactionURI))
	if err != nil {
		return err
	}
	return nil
}

// GetMessageReactionUsers 获取消息表情表态用户列表
func (o *openAPI) GetMessageReactionUsers(ctx context.Context, channelID, messageID string, emoji dto.Emoji,
	pager *dto.MessageReactionPager) (*dto.MessageReactionUsers, error) {
	if pager == nil {
		return nil, errs.ErrPagerIsNil
	}
	resp, err := o.request(ctx).
		SetPathParam("channel_id", channelID).
		SetPathParam("message_id", messageID).
		SetPathParam("emoji_type", strconv.FormatUint(uint64(emoji.Type), 10)).
		SetPathParam("emoji_id", emoji.ID).
		SetQueryParams(pager.QueryParams()).
		Get(o.getURL(messageReactionURI))
	if err != nil {
		return nil, err
	}

	messageReactionUsers := &dto.MessageReactionUsers{}
	if err := json.Unmarshal(resp.Body(), &messageReactionUsers); err != nil {
		return nil, err
	}

	return messageReactionUsers, nil
}
