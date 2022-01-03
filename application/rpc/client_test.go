/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package rpc

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
)

var rpcRunning bool

func TestAddValidator(t *testing.T) {
	validatorAddress, err := sdk.ValAddressFromBech32("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf")
	require.Nil(t, err)

	const (
		validatorName = "binance"
		rpcEndpoint   = "127.0.0.1:4040"
	)

	if !rpcRunning {
		go StartServer(rpcEndpoint)

		rpcRunning = true
	}

	database, err := db.OpenDB(constants.TestDBDir)
	require.Nil(t, err)

	defer database.Close()

	validators, err := AddValidator(db.Validator{
		Address: validatorAddress,
		Name:    validatorName,
	}, rpcEndpoint)
	require.Nil(t, err)

	validatorsGet, err := db.GetValidators()
	require.Equal(t, validators, validatorsGet)
	require.Nil(t, err)

	err = db.DeleteAllValidators()
	require.Nil(t, err)
}

func TestRemoveValidator(t *testing.T) {
	validatorAddress, err := sdk.ValAddressFromBech32("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf")
	require.Nil(t, err)

	const (
		validatorName = "binance"
		rpcEndpoint   = "127.0.0.1:4040"
	)

	if !rpcRunning {
		go StartServer(rpcEndpoint)

		rpcRunning = true
	}

	database, err := db.OpenDB(constants.TestDBDir)
	require.Nil(t, err)

	defer database.Close()

	initialValidatorsSet, err := db.GetValidators()
	require.Nil(t, err)

	validators, err := AddValidator(db.Validator{
		Address: validatorAddress,
		Name:    validatorName,
	}, rpcEndpoint)
	require.Nil(t, err)
	require.Len(t, validators, len(initialValidatorsSet)+1)

	validators, err = RemoveValidator(validatorAddress, rpcEndpoint)
	require.Equal(t, nil, err)
	require.Len(t, validators, len(initialValidatorsSet))

	validatorsGet, err2 := db.GetValidators()

	require.Equal(t, nil, err2)
	require.Equal(t, validators, validatorsGet)

	err = db.DeleteAllValidators()
	require.Nil(t, err)

	validators, err = db.GetValidators()
	require.Nil(t, err)
	require.Len(t, validators, 0)
}

func TestShowValidators(t *testing.T) {
	rpcEndpoint := "127.0.0.1:4040"

	if !rpcRunning {
		go StartServer(rpcEndpoint)

		rpcRunning = true
	}

	database, err := db.OpenDB(constants.TestDBDir)
	require.Nil(t, err)

	defer database.Close()

	validators, err := ShowValidators("", rpcEndpoint)
	require.Nil(t, err)

	validatorsGet, err2 := db.GetValidators()
	require.Equal(t, nil, err2)
	require.Equal(t, validators, validatorsGet)
}
