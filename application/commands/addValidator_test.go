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

func TestAddCommand(t *testing.T) {
	var (
		database *badger.DB
		closeFn  func()
		err      error
	)

	database, closeFn, err = test.OpenDB(t, db.OpenDB)
	defer closeFn()

	require.Nil(t, err)

	err = db.DeleteAllValidators(database)
	require.Nil(t, err)

	const (
		validatorName    = "Binance"
		validatorAddress = "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"
	)

	cmd := AddCommand()
	cmd.SetArgs([]string{validatorAddress, validatorName})

	err = cmd.Flags().Set(constants.FlagPBridgeHome, constants.TestHomeDir)
	require.Nil(t, err)

	err = cmd.Execute()
	require.Nil(t, err)

	database, closeFn, err = test.OpenDB(t, db.OpenDB)
	defer closeFn()

	require.Nil(t, err)

	address, _ := sdk.ValAddressFromBech32(validatorAddress)
	validator, err := db.GetValidator(database, address)
	require.Nil(t, err)

	require.Equal(t, validatorName, validator.Name, "Validator name does not match of that added")
	require.Equal(t, address, validator.Address, "Validator address does not match of that added")

	err = db.DeleteAllValidators(database)
	require.Nil(t, err)
}
