package apitest

import (
	"testing"

	"github.com/tencent-connect/botgo/dto"
)

func TestGetAPIPermissions(t *testing.T) {
	apiIdentify := &dto.APIPermissionDemandIdentify{}
	t.Run(
		"get api permissions", func(t *testing.T) {
			result, err := api.GetAPIPermissions(ctx, testGuildID)
			if err != nil {
				t.Error(err)
			}
			for _, v := range result.APIList {
				t.Logf("api permissions:%+v", v)
			}
			if len(result.APIList) > 0 {
				apiIdentify.Path = result.APIList[0].Path
				apiIdentify.Method = result.APIList[0].Method
			}

		},
	)
	t.Run(
		"create API permission demand", func(t *testing.T) {
			demand, err := api.RequireAPIPermissions(ctx, testGuildID, &dto.APIPermissionDemandToCreate{
				ChannelID:   testChannelID,
				APIIdentify: apiIdentify,
				Desc:        "授权链接",
			})
			if err != nil {
				t.Error(err)
			}
			t.Logf("demand:%+v", demand)
		},
	)
}
