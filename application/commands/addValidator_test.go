package commands

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/rpc"
	"os"
	"path/filepath"
	"testing"
)

func TestAddCommand(t *testing.T) {
	t.Logf("Testing AddCommand command")

	dirname, _ := os.UserHomeDir()
	db, err := db.OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	validatorName := "binance"
	validatorAddress := "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"
	var initClientCtx client.Context
	rpcEndpoint := "127.0.0.1:4040"
	validators, err := rpc.ShowValidators("", rpcEndpoint)

	cmd := AddCommand(initClientCtx)
	cmd.SetArgs([]string{validatorAddress, validatorName})
	err = cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
	validators, err = rpc.ShowValidators("", rpcEndpoint)
	length := len(validators)

	validator := validators[length-1]
	if validator.Name != validatorName {
		t.Errorf("Validaor name does not match of that added")
	}
	if validator.Address.String() != validatorAddress {
		t.Errorf("Validaor address does not match of that added")
	}

	db.Close()
}
