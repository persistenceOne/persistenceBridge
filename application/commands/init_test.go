package commands

import (
	"os"
	"testing"
)

func TestInitCommand(t *testing.T) {
	t.Logf("Testing init command")

	dirname, _ := os.UserHomeDir()
	cmd := InitCommand()
	cmd.SetArgs([]string{dirname + "/Documents/GitHub/persistenceBridge"})
	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
