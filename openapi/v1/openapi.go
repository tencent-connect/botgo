package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2" // resty 是一个优秀的 rest api 客户端，可以极大的减少开发基于 rest 标准接口求请求的封装工作量
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/errs"
	"github.com/tencent-connect/botgo/internal/log"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/version"
)

type openAPI struct {
	token   *token.Token
	timeout time.Duration
	body    interface{}
	sandbox bool
	debug   bool
}

// Setup 注册
func Setup() {
	openapi.Register(openapi.APIv1, &openAPI{})
}

// Version ...
func (o openAPI) Version() openapi.APIVersion {
	return openapi.APIv1
}

// New ...
func (o openAPI) New(token *token.Token, inSandbox bool) openapi.OpenAPI {
	return &openAPI{
		token:   token,
		timeout: 3 * time.Second,
		sandbox: inSandbox,
	}
}

func (o *openAPI) WithTimeout(duration time.Duration) openapi.OpenAPI {
	o.timeout = duration
	return o
}

func (o *openAPI) WithBody(body interface{}) openapi.OpenAPI {
	o.body = body
	return o
}

// WS ...
func (o *openAPI) WS(ctx context.Context, _ map[string]string, _ string) (*dto.WebsocketAP, error) {
	resp, err := o.request(ctx).
		SetResult(dto.WebsocketAP{}).
		Get(getURL(gatewayBotURI, o.sandbox))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.WebsocketAP), nil
}

// Me ...
func (o *openAPI) Me(ctx context.Context) (*dto.User, error) {
	resp, err := o.request(ctx).
		SetResult(dto.User{}).
		Get(getURL(userMeURI, o.sandbox))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.User), nil
}

// MeGuilds ...
func (o *openAPI) MeGuilds(ctx context.Context) ([]*dto.Guild, error) {
	resp, err := o.request(ctx).
		Get(getURL(userMeGuildsURI, o.sandbox))
	if err != nil {
		return nil, err
	}

	guilds := make([]*dto.Guild, 0)
	if err := json.Unmarshal(resp.Body(), &guilds); err != nil {
		return nil, err
	}

	return guilds, nil
}

// Message ...
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

// Messages ...
func (o *openAPI) Messages(ctx context.Context, channelID string, pager *dto.MessagesPager) ([]*dto.Message, error) {
	resp, err := o.request(ctx).
		SetPathParam("channel_id", channelID).
		SetQueryParam(string(pager.Type), pager.ID).
		SetQueryParam("limit", pager.Limit).
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

// PostMessage ...
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

// Guild ...
func (o *openAPI) Guild(ctx context.Context, guildID string) (*dto.Guild, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Guild{}).
		SetPathParam("guild_id", guildID).
		Get(getURL(guildURI, o.sandbox))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Guild), nil
}

// GuildMember ...
func (o *openAPI) GuildMember(ctx context.Context, guildID, userID string) (*dto.Member, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Member{}).
		SetPathParam("guild_id", guildID).
		SetPathParam("user_id", userID).
		Get(getURL(guildMemberURI, o.sandbox))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Member), nil
}

// GuildMembers ...
func (o *openAPI) GuildMembers(ctx context.Context,
	guildID string, pager *dto.GuildMembersPager) ([]*dto.Member, error) {
	resp, err := o.request(ctx).
		SetPathParam("guild_id", guildID).
		SetQueryParam("after", pager.After).
		SetQueryParam("limit", pager.Limit).
		Get(getURL(guildMembersURI, o.sandbox))
	if err != nil {
		return nil, err
	}

	members := make([]*dto.Member, 0)
	if err := json.Unmarshal(resp.Body(), &members); err != nil {
		return nil, err
	}

	return members, nil
}

// DeleteGuildMember ...
func (o *openAPI) DeleteGuildMember(ctx context.Context, guildID, userID string) error {
	_, err := o.request(ctx).
		SetResult(dto.Member{}).
		SetPathParam("guild_id", guildID).
		SetPathParam("user_id", userID).
		Delete(getURL(guildMemberURI, o.sandbox))
	return err
}

// Channel ...
func (o *openAPI) Channel(ctx context.Context, channelID string) (*dto.Channel, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Channel{}).
		SetPathParam("channel_id", channelID).
		Get(getURL(channelURI, o.sandbox))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Channel), nil
}

// Channels ...
func (o *openAPI) Channels(ctx context.Context, guildID string) ([]*dto.Channel, error) {
	resp, err := o.request(ctx).
		SetPathParam("guild_id", guildID).
		Get(getURL(channelsURI, o.sandbox))
	if err != nil {
		return nil, err
	}

	channels := make([]*dto.Channel, 0)
	if err := json.Unmarshal(resp.Body(), &channels); err != nil {
		return nil, err
	}

	return channels, nil
}

// PostChannel ...
func (o *openAPI) PostChannel(ctx context.Context,
	guildID string, value *dto.ChannelValueObject) (*dto.Channel, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Channel{}).
		SetPathParam("guild_id", guildID).
		SetBody(value).
		Post(getURL(channelsURI, o.sandbox))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Channel), nil
}

// PatchChannel ...
func (o *openAPI) PatchChannel(ctx context.Context,
	channelID string, value *dto.ChannelValueObject) (*dto.Channel, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Channel{}).
		SetPathParam("channel_id", channelID).
		SetBody(value).
		Patch(getURL(channelURI, o.sandbox))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Channel), nil
}

// DeleteChannel ...
func (o *openAPI) DeleteChannel(ctx context.Context, channelID string) error {
	_, err := o.request(ctx).
		SetResult(dto.Channel{}).
		SetPathParam("channel_id", channelID).
		Delete(getURL(channelURI, o.sandbox))
	return err
}

// Transport 透传请求
func (o *openAPI) Transport(ctx context.Context, method, url string, body interface{}) ([]byte, error) {
	resp, err := o.request(ctx).SetBody(body).Execute(method, url)
	return resp.Body(), err
}

func (o *openAPI) request(ctx context.Context) *resty.Request {
	client := resty.New().
		SetLogger(log.DefaultLogger).
		SetDebug(o.debug).
		SetTimeout(o.timeout).
		SetAuthToken(o.token.GetString()).
		SetAuthScheme(string(o.token.Type)).
		SetHeader("User-Agent", version.String()).
		// 设置请求之后的钩子，打印日志，判断状态码
		OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
			log.Infof("%v", respInfo(response))
			// 非成功含义的状态码，需要返回 error 供调用方识别
			if !openapi.IsSuccessStatus(response.StatusCode()) {
				return errs.New(response.StatusCode(), string(response.Body()))
			}
			return nil
		})

	return client.R().
		SetContext(ctx)
}

// respInfo 用于输出日志的时候格式化数据
func respInfo(resp *resty.Response) string {
	return fmt.Sprintf("[OPENAPI]URL:%v, Trace:%v, status:%v, body:%v",
		resp.Request.URL, resp.Header().Get(openapi.TraceIDKey), resp.Status(), string(resp.Body()))
}
