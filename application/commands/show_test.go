/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
)

func TestShowCommand(t *testing.T) {
	database, err := db.OpenDB(constants2.TestDBDir)
	require.Nil(t, err)

	err = db.DeleteAllValidators()
	require.Nil(t, err)

	valAddress, err := sdk.ValAddressFromBech32("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf")
	require.Nil(t, err)

	err = db.SetValidator(db.Validator{
		Address: valAddress,
		Name:    "binance",
	})
	require.Nil(t, err)

	database.Close()

	cmd := ShowCommand()
	err = cmd.Flags().Set(constants2.FlagPBridgeHome, constants2.TestHomeDir)
	require.Nil(t, err)
	err = cmd.Execute()
	require.Nil(t, err)

	database, err = db.OpenDB(constants2.TestDBDir)
	require.Nil(t, err)
	err = db.DeleteAllValidators()
	require.Nil(t, err)
	database.Close()
}
