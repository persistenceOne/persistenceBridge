package commands

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/rpc"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveCommand(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	database, err := db.OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	require.Nil(t, err)

	go rpc.StartServer("127.0.0.1:4040")
	validators, err := rpc.ShowValidators("", "127.0.0.1:4040")
	if len(validators) > 0 {
		if len(validators) > 0 {
			require.Nil(t, err)
			for _, validator := range validators {
				_, err = rpc.RemoveValidator(validator.Address, "127.0.0.1:4040")
				require.Nil(t, err)
			}
		}
	}
	validatorName1 := "Binance"
	validatorAddress1 := "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"
	validatorName2 := "myAddress"
	validatorAddress2 := "cosmosvaloper18r9630ruvesw76h2qand6pvdzjctpp6q4dlgm5"

	valAddress1, err := sdk.ValAddressFromBech32(validatorAddress1)
	require.Nil(t, err)
	valAddress2, err := sdk.ValAddressFromBech32(validatorAddress2)
	require.Nil(t, err)

	_, err = rpc.AddValidator(db.Validator{
		Address: valAddress1,
		Name:    validatorName1,
	}, "127.0.0.1:4040")
	require.Nil(t, err)

	_, err = rpc.AddValidator(db.Validator{
		Address: valAddress2,
		Name:    validatorName2,
	}, "127.0.0.1:4040")
	require.Nil(t, err)

	cmd := RemoveCommand()
	cmd.SetArgs([]string{validatorAddress2})
	err = cmd.Execute()
	require.Nil(t, err)

	database.Close()
	database, err = db.OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	require.Nil(t, err)
	validator1, err := db.GetValidator(valAddress1)
	require.Nil(t, err)

	require.Equal(t, validatorName1, validator1.Name, "Validaor name does not match of that added")
	require.Equal(t, valAddress1.String(), validator1.Address.String(), "Validator address does not match of that added")

	cmd = RemoveCommand()
	cmd.SetArgs([]string{validatorAddress1})
	err = cmd.Execute()

	require.Equal(t, "need to have at least one validator to redelegate to", err.Error())

	database.Close()
}
