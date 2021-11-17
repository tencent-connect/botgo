package v1

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tencent-connect/botgo/dto"
)

// ChannelPermissions 获取指定子频道的权限
func (o *openAPI) ChannelPermissions(ctx context.Context, channelID, userID string) (*dto.ChannelPermissions, error) {
	rsp, err := o.request(ctx).
		SetResult(dto.ChannelPermissions{}).
		SetPathParam("channel_id", channelID).
		SetPathParam("user_id", userID).
		Get(getURL(channelPermissionsURI, o.sandbox))
	if err != nil {
		return nil, err
	}
	return rsp.Result().(*dto.ChannelPermissions), nil
}

// PutChannelPermissions 修改指定子频道的权限
func (o *openAPI) PutChannelPermissions(ctx context.Context, channelID, userID string,
	p *dto.UpdateChannelPermissions) error {
	if p.Add != "" {
		if _, err := strconv.ParseUint(p.Add, 10, 64); err != nil {
			return fmt.Errorf("invalid parameter add: %v", err)
		}
	}
	if p.Remove != "" {
		if _, err := strconv.ParseUint(p.Remove, 10, 64); err != nil {
			return fmt.Errorf("invalid parameter remove: %v", err)
		}
	}
	_, err := o.request(ctx).
		SetPathParam("channel_id", channelID).
		SetPathParam("user_id", userID).
		SetBody(p).
		Put(getURL(channelPermissionsURI, o.sandbox))
	return err
}
