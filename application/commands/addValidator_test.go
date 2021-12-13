/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddCommand(t *testing.T) {
	database, err := db.OpenDB(constants2.TestDbDir)
	require.Nil(t, err)
	err = db.DeleteAllValidators()
	require.Nil(t, err)
	database.Close()

	validatorName := "Binance"
	validatorAddress := "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

	cmd := AddCommand()
	cmd.SetArgs([]string{validatorAddress, validatorName})
	err = cmd.Flags().Set(constants2.FlagPBridgeHome, constants2.TestHomeDir)
	require.Equal(t, nil, err)
	err = cmd.Execute()
	require.Nil(t, err)

	database, err = db.OpenDB(constants2.TestDbDir)
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
