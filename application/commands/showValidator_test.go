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

func TestShowCommand(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	database, err := db.OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	require.Nil(t, err)

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

	validatorName := "binance"
	validatorAddress := "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"
	valAddress, err := sdk.ValAddressFromBech32(validatorAddress)

	cmd := ShowCommand()
	err = cmd.Execute()
	require.Equal(t, nil, err)

	validators, err = rpc.AddValidator(db.Validator{
		Address: valAddress,
		Name:    validatorName,
	}, "127.0.0.1:4040")

	cmd = ShowCommand()
	err = cmd.Execute()
	require.Equal(t, nil, err)
	database.Close()
}
