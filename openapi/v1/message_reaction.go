package v1

import (
	"context"
	"strconv"

	"github.com/tencent-connect/botgo/dto"
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
