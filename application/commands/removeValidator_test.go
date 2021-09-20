package commands

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRemoveCommand(t *testing.T) {
	database, err := db.OpenDB(constants2.TestDbDir)
	require.Nil(t, err)
	err = db.DeleteAllValidators()
	require.Nil(t, err)

	validatorName1 := "Binance"
	validatorAddress1 := "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"
	validatorName2 := "myAddress"
	validatorAddress2 := "cosmosvaloper18r9630ruvesw76h2qand6pvdzjctpp6q4dlgm5"

	valAddress1, err := sdk.ValAddressFromBech32(validatorAddress1)
	require.Nil(t, err)
	valAddress2, err := sdk.ValAddressFromBech32(validatorAddress2)
	require.Nil(t, err)

	validator1 := db.Validator{
		Address: valAddress1,
		Name:    validatorName1,
	}
	err = db.SetValidator(validator1)
	require.Nil(t, err)

	err = db.SetValidator(db.Validator{
		Address: valAddress2,
		Name:    validatorName2,
	})
	require.Nil(t, err)

	err = database.Close()
	require.Nil(t, err)

	cmd := RemoveCommand()
	err = cmd.Flags().Set(constants2.FlagPBridgeHome, constants2.TestHomeDir)
	require.Nil(t, err)
	cmd.SetArgs([]string{validatorAddress2})
	err = cmd.Execute()
	require.Nil(t, err)

	database, err = db.OpenDB(constants2.TestDbDir)
	require.Nil(t, err)

	validators, err := db.GetValidators()
	require.Nil(t, err)

	require.Equal(t, []db.Validator{validator1}, validators)

	err = db.DeleteAllValidators()
	require.Nil(t, err)
	database.Close()
}
