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

func TestRemoveCommand(t *testing.T) {
	t.Logf("Testing RemoveCommand command")

	dirname, _ := os.UserHomeDir()
	database, err := db.OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	validatorName1 := "binance"
	validatorAddress1 := "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"
	validatorName2 := "my_address"
	validatorAddress2 := "cosmosvaloper18r9630ruvesw76h2qand6pvdzjctpp6q4dlgm5"
	rpcEndpoint := "127.0.0.1:4040"
	var initClientCtx client.Context
	validators, err := rpc.ShowValidators("", rpcEndpoint)

	valAddress1, err := sdk.ValAddressFromBech32(validatorAddress1)
	valAddress2, err := sdk.ValAddressFromBech32(validatorAddress2)

	for _, validator := range validators {
		_, _ = rpc.RemoveValidator(validator.Address, rpcEndpoint)
	}

	validators, err = rpc.AddValidator(db.Validator{
		Address: valAddress1,
		Name:    validatorName1,
	}, rpcEndpoint)

	validators, err = rpc.AddValidator(db.Validator{
		Address: valAddress2,
		Name:    validatorName2,
	}, rpcEndpoint)

	cmd := RemoveCommand(initClientCtx)
	cmd.SetArgs([]string{validatorAddress2})
	err = cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	validators, err = rpc.ShowValidators("", rpcEndpoint)

	if validators[0].Name != validatorName1 && validators[0].Address.String() != validatorAddress1 {
		t.Fatalf("expected Validator name %s got %s and Validator Address %s got %s \n", validatorName1, validators[0].Name, validatorAddress1, validators[0].Address.String())
	}

	cmd = RemoveCommand(initClientCtx)
	cmd.SetArgs([]string{validatorAddress1})
	err = cmd.Execute()

	if err.Error() != "need to have at least one validator to redelegate to" {
		t.Errorf("Expected an error due to no validator to redelgate")
	}

	database.Close()
}
