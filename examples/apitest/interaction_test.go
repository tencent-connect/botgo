package apitest

import (
	"encoding/json"
	"testing"

	"github.com/tencent-connect/botgo/dto"
)

func TestInteractions(t *testing.T) {
	t.Run(
		"put interaction", func(t *testing.T) {
			body, _ := json.Marshal(
				dto.InteractionData{
					Name:     "interaction",
					Type:     2,
					Resolved: json.RawMessage("test"),
				},
			)
			err := api.PutInteraction(ctx, testInteractionD, string(body))
			if err != nil {
				t.Error(err)
			}
		},
	)
}
