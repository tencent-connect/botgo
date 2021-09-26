package websocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tencent-connect/botgo/dto"
)

func TestRegisterHandlers(t *testing.T) {
	var guild GuildEventHandler = func(event *dto.WSPayload, data *dto.WSGuildData) error {
		return nil
	}
	var message MessageEventHandler = func(event *dto.WSPayload, data *dto.WSMessageData) error {
		return nil
	}

	t.Run("test intent", func(t *testing.T) {
		i := RegisterHandlers(guild, message)
		assert.Equal(t, dto.IntentGuildMessages, i&dto.IntentGuildMessages)
		assert.Equal(t, dto.IntentGuilds, i&dto.IntentGuilds)
	})
}
