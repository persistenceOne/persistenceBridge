package commands

import (
	"github.com/cosmos/cosmos-sdk/client"
	"os"
	"testing"
)

func TestStartCommand(t *testing.T) {
	t.Logf("Testing start command")

	dirname, _ := os.UserHomeDir()
	var initClientCtx client.Context
	cmd := StartCommand(initClientCtx)
	cmd.SetArgs([]string{dirname + "/Documents/GitHub/persistenceBridge"})
	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
