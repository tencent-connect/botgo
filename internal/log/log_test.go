package log

import (
	"testing"
)

func TestDebug(t *testing.T) {
	Debug("debug log")
	Error("error log")
	Warn("warn log")
	Info("info log")
	Debugf("%s log", "debugf")
	Errorf("%s log", "errorf")
	Warnf("%s log", "warnf")
	Infof("%s log", "infof")
}
