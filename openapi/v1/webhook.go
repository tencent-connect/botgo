package v1

import (
	"context"
	"encoding/json"

	"github.com/tencent-connect/botgo/dto"
)

// CreateSession 创建一个新的 http 事件回调
func (o *openAPI) CreateSession(ctx context.Context, identity dto.HTTPIdentity) (*dto.HTTPReady, error) {
	resp, err := o.request(ctx).
		SetResult(dto.HTTPReady{}).
		SetBody(identity).
		Post(o.getURL(httpSessionsURI))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.HTTPReady), nil
}

// CheckSessions 定期检查 http 回调 session 的健康情况，服务端会自动 resume 非活跃状态的 session
func (o *openAPI) CheckSessions(ctx context.Context) ([]*dto.HTTPSession, error) {
	resp, err := o.request(ctx).
		SetQueryParam("action", "check").
		Patch(o.getURL(httpSessionsURI))
	if err != nil {
		return nil, err
	}

	sessions := make([]*dto.HTTPSession, 0)
	if err := json.Unmarshal(resp.Body(), &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

// GetActiveSessionList 拉取活跃的 http session 列表
func (o *openAPI) SessionList(ctx context.Context) ([]*dto.HTTPSession, error) {
	resp, err := o.request(ctx).
		Get(o.getURL(httpSessionsURI))
	if err != nil {
		return nil, err
	}

	sessions := make([]*dto.HTTPSession, 0)
	if err := json.Unmarshal(resp.Body(), &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

// RemoveSession 停止某个 session
func (o *openAPI) RemoveSession(ctx context.Context, sessionID string) error {
	_, err := o.request(ctx).
		SetPathParam("session_id", sessionID).
		Delete(o.getURL(httpSessionURI))
	if err != nil {
		return err
	}

	return nil
}
