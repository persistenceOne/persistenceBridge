package utilities

import "testing"

func TestLogs(t *testing.T) {
	Token = "1805701775:AAHu-P76O2oaTDGXiBmbylpARaBr2Pdcr8s"
	ChatId = -577933520
	LogInfo("test")
	LogWarning("test")
	LogFatal("test")
	LogDebug("test")
}
