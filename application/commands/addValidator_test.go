package commands

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestAddCommand(t *testing.T) {
	t.Logf("Testing AddCommand command")

	dirname, _ := os.UserHomeDir()

	validatorName := "Binance"
	validatorAddress := "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

	cmd := AddCommand()
	cmd.SetArgs([]string{validatorAddress, validatorName})
	err := cmd.Execute()
	require.Nil(t, err)

	database, err := db.OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	require.Nil(t, err)
	address, _ := sdk.ValAddressFromBech32(validatorAddress)
	validators, err := db.GetValidator(address)
	require.Nil(t, err)

	require.Equal(t, validatorName, validators.Name, "Validaor name does not match of that added")
	require.Equal(t, address, validators.Address,  "Validator address does not match of that added")

	database.Close()
}
