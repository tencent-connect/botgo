package apitest

import (
	"testing"
)

func TestRetractDMMessage(t *testing.T) {
	msgID := "10c7fac2c28dcac27a1a1231343431313532313831383136323933383420801e28003095a8683824402448cefba48e0650b1acf8fa05"
	t.Run(
		"私信消息撤回", func(t *testing.T) {
			err := api.RetractDMMessage(ctx, "6234704349443091672", msgID)
			if err != nil {
				t.Error(err)
			}
			t.Logf("msg id : %v, is deleted", msgID)
		},
	)
}
