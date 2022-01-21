package v1

import (
	"context"

	"github.com/tencent-connect/botgo/dto"
)

// WS 获取带分片 WSS 接入点
func (o *openAPI) WS(ctx context.Context, _ map[string]string, _ string) (*dto.WebsocketAP, error) {
	resp, err := o.request(ctx).
		SetResult(dto.WebsocketAP{}).
		Get(o.getURL(gatewayBotURI))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.WebsocketAP), nil
}
