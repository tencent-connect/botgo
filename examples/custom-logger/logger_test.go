package main

import (
	"testing"
)

func TestFileLogger_Debug(t *testing.T) {
	l, err := New("./", DebugLevel)
	t.Error(err)
	l.Debug("abc")
}
