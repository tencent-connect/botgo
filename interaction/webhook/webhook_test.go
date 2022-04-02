package webhook

import (
	"strings"
	"testing"
)

func TestGenHeartbeatACK(t *testing.T) {
	j := GenHeartbeatACK(1314)
	if !strings.Contains(j, "11") || !strings.Contains(j, "1314") {
		t.Error("GenHeartbeatACK error")
	}
}

func TestGenDispatchACK(t *testing.T) {
	j := GenDispatchACK(false)
	if !strings.Contains(j, "12") || !strings.Contains(j, "1") {
		t.Error("GenDispatchACK error")
	}
	j = GenDispatchACK(true)
	if !strings.Contains(j, "12") || !strings.Contains(j, "0") {
		t.Error("GenDispatchACK error")
	}
}
