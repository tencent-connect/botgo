package apitest

import (
	"testing"
	"time"
)

func TestPins(t *testing.T) {

	t.Run(
		"add pins", func(t *testing.T) {
			pins, err := api.AddPins(
				ctx, testChannelID, testMessageID,
			)
			if err != nil {
				t.Error(err)
			}
			t.Logf("pins:%+v", pins)
		},
	)
	t.Run(
		"get pins", func(t *testing.T) {
			time.Sleep(3 * time.Second)
			pins, err := api.GetPins(ctx, testChannelID)
			if err != nil {
				t.Error(err)
			}
			t.Logf("pins:%+v", pins)

		},
	)
	t.Run(
		"delete pins", func(t *testing.T) {
			time.Sleep(3 * time.Second)
			err := api.DeletePins(ctx, testChannelID, testMessageID)
			if err != nil {
				t.Error(err)
			}
		},
	)
	t.Run(
		"clean pins no check message id", func(t *testing.T) {
			time.Sleep(3 * time.Second)
			err := api.CleanPins(ctx, testChannelID)
			if err != nil {
				t.Error(err)
			}
		},
	)
}
