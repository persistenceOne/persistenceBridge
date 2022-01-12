//go:build integration

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestShowCommand(t *testing.T) {
	{
		database, closeFn, err := test.OpenDB(t, db.OpenDB)
		defer closeFn()

		require.Nil(t, err)

		err = db.DeleteAllValidators(database)
		require.Nil(t, err)

		valAddress, err := sdk.ValAddressFromBech32("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf")
		require.Nil(t, err)

		err = db.SetValidator(database, db.Validator{
			Address: valAddress,
			Name:    "binance",
		})
		require.Nil(t, err)
	}

	{
		cmd := ShowCommand()
		err := cmd.Flags().Set(constants.FlagPBridgeHome, constants.TestHomeDir())
		require.Nil(t, err)

		err = cmd.Execute()
		require.Nil(t, err)

		var (
			database *badger.DB
			closeFn  func()
		)

		database, closeFn, err = test.OpenDB(t, db.OpenDB)
		defer closeFn()

		require.Nil(t, err)

		err = db.DeleteAllValidators(database)
		require.Nil(t, err)
	}
}
