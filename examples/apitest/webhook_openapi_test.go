package apitest

import (
	"testing"

	"github.com/tencent-connect/botgo/dto"
)

func TestWebhook(t *testing.T) {
	var sessionID string
	t.Run(
		"会话鉴权", func(t *testing.T) {
			rsp, err := api.CreateSession(
				ctx, dto.HTTPIdentity{
					Callback: "https://demo.app.tcloudbase.com/callback",
					Intents:  1024, Shards: [2]uint32{1, 4},
				},
			)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%v", rsp)
			sessionID = rsp.SessionID
		},
	)
	t.Run(
		"拉取会话", func(t *testing.T) {
			rsp, err := api.SessionList(ctx)
			if err != nil {
				t.Error(err)
			}
			for _, v := range rsp {
				t.Logf("%v", v)
			}
		},
	)
	t.Run(
		"检测会话", func(t *testing.T) {
			rsp, err := api.CheckSessions(ctx)
			if err != nil {
				t.Error(err)
			}
			for _, v := range rsp {
				if v.State != "ACTIVE" {
					t.Errorf("data status error, data: %v", v)
				} else {
					t.Logf("%v", v)
				}
			}
		},
	)
	t.Run(
		"关闭会话", func(t *testing.T) {
			err := api.RemoveSession(ctx, sessionID)
			if err != nil {
				t.Error(err)
			} else {
				t.Logf("stop success")
			}
		},
	)
}
