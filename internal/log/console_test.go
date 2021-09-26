package log

import (
	"testing"
)

func Test_log(t *testing.T) {
	t.Run("output", func(t *testing.T) {
		output("Info", "abc", "def")
	})

	t.Run("Debug", func(t *testing.T) {
		c := consoleLogger{}
		c.Debug("debug log")
	})

	t.Run("Debugf", func(t *testing.T) {
		c := consoleLogger{}
		c.Debugf("debugf %s", "log")
	})
}
