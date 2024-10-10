package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/errs"
	"github.com/tencent-connect/botgo/openapi/options"
)

// Message 拉取单条消息
func (o *openAPI) Message(ctx context.Context, channelID string, messageID string, opt ...options.Option) (
	*dto.Message, error) {
	reqCMD := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("channel_id", channelID).
		SetPathParam("message_id", messageID)

	resp, err := baseRequest(ctx, reqCMD, http.MethodGet, o.getURL(messagesURI), opt...)
	if err != nil {
		return nil, err
	}

	// 兼容处理
	result := resp.Result().(*dto.Message)
	if result.ID == "" {
		body := gjson.Get(resp.String(), "message")
		if err := json.Unmarshal([]byte(body.String()), result); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// Messages 拉取消息列表
func (o *openAPI) Messages(ctx context.Context, channelID string, pager *dto.MessagesPager, opt ...options.Option) (
	[]*dto.Message, error) {
	if pager == nil {
		return nil, errs.ErrPagerIsNil
	}
	reqCMD := o.request(ctx).
		SetPathParam("channel_id", channelID).
		SetQueryParams(pager.QueryParams())

	resp, err := baseRequest(ctx, reqCMD, http.MethodGet, o.getURL(messagesURI), opt...)
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
func (o *openAPI) PostMessage(ctx context.Context, channelID string, msg *dto.MessageToCreate,
	opt ...options.Option) (*dto.Message, error) {
	reqCMD := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("channel_id", channelID).
		SetBody(msg)

	resp, err := baseRequest(ctx, reqCMD, http.MethodPost, o.getURL(messagesURI), opt...)
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.Message), nil
}

// PatchMessage 编辑消息
func (o *openAPI) PatchMessage(ctx context.Context,
	channelID string, messageID string, msg *dto.MessageToCreate, opt ...options.Option) (*dto.Message, error) {
	reqCMD := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("channel_id", channelID).
		SetPathParam("message_id", messageID).
		SetBody(msg)

	resp, err := baseRequest(ctx, reqCMD, http.MethodPatch, o.getURL(messageURI), opt...)
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Message), nil
}

// RetractMessage 撤回消息
func (o *openAPI) RetractMessage(ctx context.Context,
	channelID, msgID string, opt ...options.Option) error {
	reqCMD := o.request(ctx).
		SetPathParam("channel_id", channelID).
		SetPathParam("message_id", string(msgID))

	_, err := baseRequest(ctx, reqCMD, http.MethodDelete, o.getURL(messageURI), opt...)
	return err
}

// RetractC2CMessage 撤回C2C消息
func (o *openAPI) RetractC2CMessage(ctx context.Context,
	userID, msgID string, opt ...options.Option) error {
	reqCMD := o.request(ctx).
		SetPathParam("user_id", userID).
		SetPathParam("message_id", msgID)
	_, err := baseRequest(ctx, reqCMD, http.MethodDelete, o.getURL(retractC2cMessageURI), opt...)
	return err
}

// RetractGroupMessage 撤回群消息
func (o *openAPI) RetractGroupMessage(ctx context.Context,
	groupID, msgID string, opt ...options.Option) error {
	reqCMD := o.request(ctx).
		SetPathParam("group_id", groupID).
		SetPathParam("message_id", msgID)

	_, err := baseRequest(ctx, reqCMD, http.MethodDelete, o.getURL(retractGroupMessageURI), opt...)
	return err
}

// PostSettingGuide 发送设置引导消息, atUserID为要at的用户
func (o *openAPI) PostSettingGuide(ctx context.Context,
	channelID string, atUserIDs []string, opt ...options.Option) (*dto.Message, error) {
	var content string
	for _, userID := range atUserIDs {
		content += fmt.Sprintf("<@%s>", userID)
	}
	msg := &dto.SettingGuideToCreate{
		Content: content,
	}
	reqCMD := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("channel_id", channelID).
		SetBody(msg)

	resp, err := baseRequest(ctx, reqCMD, http.MethodPost, o.getURL(settingGuideURI), opt...)
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.Message), nil
}

func getGroupURIBySendType(msgType dto.SendType) uri {
	switch msgType {
	case dto.RichMedia:
		return groupRichMediaURI
	default:
		return groupMessagesURI
	}
}

// PostGroupMessage 回复群消息
func (o *openAPI) PostGroupMessage(ctx context.Context, groupID string, msg dto.APIMessage,
	opt ...options.Option) (*dto.Message, error) {
	reqCMD := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("group_id", groupID).
		SetBody(msg)

	resp, err := baseRequest(ctx, reqCMD, http.MethodPost, o.getURL(getGroupURIBySendType(msg.GetSendType())), opt...)
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.Message), nil
}

func getC2CURIBySendType(msgType dto.SendType) uri {
	switch msgType {
	case dto.RichMedia:
		return c2cRichMediaURI
	default:
		return c2cMessagesURI
	}
}

// PostC2CMessage 回复C2C消息
func (o *openAPI) PostC2CMessage(ctx context.Context, userID string, msg dto.APIMessage,
	opt ...options.Option) (*dto.Message, error) {
	reqCMD := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("user_id", userID).
		SetBody(msg)

	resp, err := baseRequest(ctx, reqCMD, http.MethodPost, o.getURL(getC2CURIBySendType(msg.GetSendType())), opt...)
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.Message), nil
}

func baseRequest(ctx context.Context, reqCMD *resty.Request, method, url string, opt ...options.Option) (
	*resty.Response, error) {
	opts := getOptions(ctx, opt...)
	if opts.URL != "" {
		url = opts.URL
	}
	if opts.HideTip {
		reqCMD = reqCMD.SetQueryParam("hidetip", "true")
	}

	return reqCMD.Execute(method, url)
}

func getOptions(_ context.Context, opt ...options.Option) *options.Options {
	opts := &options.Options{}
	for _, o := range opt {
		o(opts)
	}
	return opts
}
