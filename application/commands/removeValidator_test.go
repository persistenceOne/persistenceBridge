//go:build integration

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestRemoveCommand(t *testing.T) {
	const (
		validatorName1    = "Binance"
		validatorAddress1 = "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"
		validatorName2    = "myAddress"
		validatorAddress2 = "cosmosvaloper18r9630ruvesw76h2qand6pvdzjctpp6q4dlgm5"
	)

	valAddress1, err := sdk.ValAddressFromBech32(validatorAddress1)
	require.Nil(t, err)

	valAddress2, err := sdk.ValAddressFromBech32(validatorAddress2)
	require.Nil(t, err)

	validator1 := db.Validator{
		Address: valAddress1,
		Name:    validatorName1,
	}

	{
		database, closeFn, err := test.OpenDB(t, db.OpenDB)
		defer closeFn()

		require.Nil(t, err)

		err = db.DeleteAllValidators(database)
		require.Nil(t, err)

		err = db.SetValidator(database, validator1)
		require.Nil(t, err)

		err = db.SetValidator(database, db.Validator{
			Address: valAddress2,
			Name:    validatorName2,
		})
		require.Nil(t, err)

		err = database.Close()
		require.Nil(t, err)

		cmd := RemoveCommand()

		err = cmd.Flags().Set(constants.FlagPBridgeHome, constants.TestHomeDir)
		require.Nil(t, err)

		cmd.SetArgs([]string{validatorAddress2})

		err = cmd.Execute()
		require.Nil(t, err)
	}

	{
		database, closeFn, err := test.OpenDB(t, db.OpenDB)
		defer closeFn()

		require.Nil(t, err)

		validators, err := db.GetValidators(database)
		require.Nil(t, err)

		require.Equal(t, []db.Validator{validator1}, validators)

		err = db.DeleteAllValidators(database)
		require.Nil(t, err)
	}
}
