package logger

import "testing"

func TestDebug(t *testing.T) {
	Debug(AllFormat, "test", "request", "response", nil)
}
