package v1

import (
	"context"
	"encoding/json"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/errs"
)

// Message 拉取单条消息
func (o *openAPI) Message(ctx context.Context, channelID string, messageID string) (*dto.Message, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("channel_id", channelID).
		SetPathParam("message_id", messageID).
		Get(getURL(messageURI, o.sandbox))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Message), nil
}

// Messages 拉取消息列表
func (o *openAPI) Messages(ctx context.Context, channelID string, pager *dto.MessagesPager) ([]*dto.Message, error) {
	if pager == nil {
		return nil, errs.ErrPagerIsNil
	}
	resp, err := o.request(ctx).
		SetPathParam("channel_id", channelID).
		SetQueryParams(pager.QueryParams()).
		Get(getURL(messagesURI, o.sandbox))
	if err != nil {
		return nil, err
	}

	messages := make([]*dto.Message, 0)
	if err := json.Unmarshal(resp.Body(), &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// PostMessage 发消息
func (o *openAPI) PostMessage(ctx context.Context, channelID string, msg *dto.MessageToCreate) (*dto.Message, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("channel_id", channelID).
		SetBody(msg).
		Post(getURL(messagesURI, o.sandbox))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Message), nil
}

// RetractMessage 撤回消息
func (o *openAPI) RetractMessage(ctx context.Context, channelID, msgID string) error {
	_, err := o.request(ctx).
		SetPathParam("channel_id", channelID).
		SetPathParam("message_id", string(msgID)).
		Delete(getURL(messageURI, o.sandbox))
	return err
}
