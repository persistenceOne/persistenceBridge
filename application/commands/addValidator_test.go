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
)

func TestAddCommand(t *testing.T) {
	database, err := db.OpenDB(constants.TestDBDir)
	require.Nil(t, err)

	err = db.DeleteAllValidators()
	require.Nil(t, err)

	database.Close()

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

	database, err = db.OpenDB(constants.TestDBDir)
	require.Nil(t, err)

	address, _ := sdk.ValAddressFromBech32(validatorAddress)
	validator, err := db.GetValidator(address)
	require.Nil(t, err)

	require.Equal(t, validatorName, validator.Name, "Validator name does not match of that added")
	require.Equal(t, address, validator.Address, "Validator address does not match of that added")

	err = db.DeleteAllValidators()
	require.Nil(t, err)
	database.Close()
}
