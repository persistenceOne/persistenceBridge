package commands

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/rpc"
	"os"
	"path/filepath"
	"testing"
)

func TestShowCommand(t *testing.T) {
	t.Logf("Testing ShowCommand command")

	dirname, _ := os.UserHomeDir()
	database, err := db.OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	validatorName := "binance"
	validatorAddress := "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"
	rpcEndpoint := "127.0.0.1:4040"
	var initClientCtx client.Context
	validators, err := rpc.ShowValidators("", rpcEndpoint)

	valAddress, err := sdk.ValAddressFromBech32(validatorAddress)

	for _, validator := range validators {
		_, _ = rpc.RemoveValidator(validator.Address, rpcEndpoint)
	}

	cmd := ShowCommand(initClientCtx)
	cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	validators, err = rpc.AddValidator(db.Validator{
		Address: valAddress,
		Name:    validatorName,
	}, rpcEndpoint)

	cmd = ShowCommand(initClientCtx)
	if err != nil {
		t.Fatal(err)
	}
	database.Close()
}
