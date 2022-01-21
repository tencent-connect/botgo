package v1

import (
	"context"
	"encoding/json"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/errs"
)

// Me 拉取当前用户的信息
func (o *openAPI) Me(ctx context.Context) (*dto.User, error) {
	resp, err := o.request(ctx).
		SetResult(dto.User{}).
		Get(o.getURL(userMeURI))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.User), nil
}

// MeGuilds 拉取当前用户加入的频道列表
func (o *openAPI) MeGuilds(ctx context.Context, pager *dto.GuildPager) ([]*dto.Guild, error) {
	if pager == nil {
		return nil, errs.ErrPagerIsNil
	}
	resp, err := o.request(ctx).
		SetQueryParams(pager.QueryParams()).
		Get(o.getURL(userMeGuildsURI))
	if err != nil {
		return nil, err
	}

	guilds := make([]*dto.Guild, 0)
	if err := json.Unmarshal(resp.Body(), &guilds); err != nil {
		return nil, err
	}

	return guilds, nil
}
