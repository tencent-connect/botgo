package apitest

import (
	"testing"
)

func TestMessageSetting(t *testing.T) {
	t.Run(
		"get message setting", func(t *testing.T) {
			settingInfo, err := api.GetMessageSetting(
				ctx, testGuildID,
			)
			if err != nil {
				t.Error(err)
			}
			t.Logf("settingIfno:%+v", settingInfo)
		},
	)
}
