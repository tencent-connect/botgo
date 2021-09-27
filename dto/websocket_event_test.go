package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_transposeIntentEventMap(t *testing.T) {
	t.Run("transpose", func(t *testing.T) {
		re := transposeIntentEventMap(intentEventMap)
		assert.Equal(t, re[EventAudioFinish], IntentAudio)
		assert.Equal(t, re[EventAudioOffMic], IntentAudio)
		assert.Equal(t, re[EventChannelCreate], IntentGuilds)
	})
}
