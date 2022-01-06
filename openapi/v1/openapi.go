// Package v1 是 openapi v1 版本的实现。
package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2" // resty 是一个优秀的 rest api 客户端，可以极大的减少开发基于 rest 标准接口求请求的封装工作量
	"github.com/tencent-connect/botgo/errs"
	"github.com/tencent-connect/botgo/log"
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
	trace   string // trace id
}

// Setup 注册
func Setup() {
	openapi.Register(openapi.APIv1, &openAPI{})
}

// Version 创建当前版本
func (o *openAPI) Version() openapi.APIVersion {
	return openapi.APIv1
}

// TraceID 获取 trace id
func (o *openAPI) TraceID() string {
	return o.trace
}

// New 生成一个实例
func (o *openAPI) New(token *token.Token, inSandbox bool) openapi.OpenAPI {
	return &openAPI{
		token:   token,
		timeout: 3 * time.Second,
		sandbox: inSandbox,
	}
}

// WithTimeout 设置请求接口超时时间
func (o *openAPI) WithTimeout(duration time.Duration) openapi.OpenAPI {
	o.timeout = duration
	return o
}

// WithBody 设置 body，如果 openapi 提供设置 body 的功能，则需要自行识别 body 类型
func (o *openAPI) WithBody(body interface{}) openapi.OpenAPI {
	o.body = body
	return o
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
		SetPreRequestHook(func(client *resty.Client, request *http.Request) error {
			// 执行请求前过滤器
			// 由于在 `OnBeforeRequest` 的时候，request 还没生成，所以 filter 不能使用，所以放到 `PreRequestHook`
			return openapi.DoReqFilterChains(request, nil)
		}).
		// 设置请求之后的钩子，打印日志，判断状态码
		OnAfterResponse(func(client *resty.Client, resp *resty.Response) error {
			log.Infof("%v", respInfo(resp))
			// 执行请求后过滤器
			if err := openapi.DoRespFilterChains(resp.Request.RawRequest, resp.RawResponse); err != nil {
				return err
			}
			o.trace = resp.Header().Get(openapi.TraceIDKey)
			// 非成功含义的状态码，需要返回 error 供调用方识别
			if !openapi.IsSuccessStatus(resp.StatusCode()) {
				return errs.New(resp.StatusCode(), string(resp.Body()), o.trace)
			}
			return nil
		})

	return client.R().
		SetContext(ctx)
}

// respInfo 用于输出日志的时候格式化数据
func respInfo(resp *resty.Response) string {
	bodyJSON, _ := json.Marshal(resp.Request.Body)
	return fmt.Sprintf("[OPENAPI]%v %v, trace:%v, status:%v, elapsed:%v req: %v, resp: %v",
		resp.Request.Method,
		resp.Request.URL,
		resp.Header().Get(openapi.TraceIDKey),
		resp.Status(),
		resp.Time(),
		string(bodyJSON),
		string(resp.Body()),
	)
}
