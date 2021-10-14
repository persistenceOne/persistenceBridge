package commands

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShowCommand(t *testing.T) {
	test.SetTestConfig()
	database, err := db.OpenDB(constants2.TestDbDir)
	require.Nil(t, err)
	err = db.DeleteAllValidators()
	require.Nil(t, err)

	valAddress, err := sdk.ValAddressFromBech32("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf")
	err = db.SetValidator(db.Validator{
		Address: valAddress,
		Name:    "binance",
	})
	require.Equal(t, nil, err)
	database.Close()

	cmd := ShowCommand()
	err = cmd.Flags().Set(constants2.FlagPBridgeHome, constants2.TestHomeDir)
	require.Equal(t, nil, err)
	err = cmd.Execute()
	require.Equal(t, nil, err)

	database, err = db.OpenDB(constants2.TestDbDir)
	require.Nil(t, err)
	err = db.DeleteAllValidators()
	require.Nil(t, err)
	database.Close()
}
