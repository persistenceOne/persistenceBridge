package rpc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/stretchr/testify/require"
	"testing"
)

var rpcRunning bool

func TestAddValidator(t *testing.T) {
	validatorAddress, err := sdk.ValAddressFromBech32("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf")
	require.Equal(t, nil, err)

	validatorName := "binance"
	rpcEndpoint := "127.0.0.1:4040"
	if !rpcRunning {
		go StartServer(rpcEndpoint)
		rpcRunning = true
	}

	database, err := db.OpenDB(constants2.TestDbDir)
	require.Nil(t, err)
	defer database.Close()

	validators, err := AddValidator(db.Validator{
		Address: validatorAddress,
		Name:    validatorName,
	}, rpcEndpoint)

	validatorsGet, err := db.GetValidators()
	require.Equal(t, validators, validatorsGet)
	require.Equal(t, nil, err)

	err = db.DeleteAllValidators()
	require.Nil(t, err)
}

func TestRemoveValidator(t *testing.T) {
	validatorAddress, err := sdk.ValAddressFromBech32("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf")
	require.Equal(t, nil, err)

	validatorName := "binance"
	rpcEndpoint := "127.0.0.1:4040"
	if !rpcRunning {
		go StartServer(rpcEndpoint)
		rpcRunning = true
	}

	database, err := db.OpenDB(constants2.TestDbDir)
	require.Nil(t, err)
	defer database.Close()

	validators, err2 := AddValidator(db.Validator{
		Address: validatorAddress,
		Name:    validatorName,
	}, rpcEndpoint)

	validators, err2 = RemoveValidator(validatorAddress, rpcEndpoint)
	require.Equal(t, nil, err2)

	validatorsGet, err2 := db.GetValidators()

	require.Equal(t, nil, err2)
	require.Equal(t, validators, validatorsGet)
	err = db.DeleteAllValidators()
	require.Nil(t, err)
}

func TestShowValidators(t *testing.T) {
	rpcEndpoint := "127.0.0.1:4040"
	if !rpcRunning {
		go StartServer(rpcEndpoint)
		rpcRunning = true
	}

	database, err := db.OpenDB(constants2.TestDbDir)
	require.Nil(t, err)
	defer database.Close()

	validators, err := ShowValidators("", rpcEndpoint)
	require.Equal(t, nil, err)

	validatorsGet, err2 := db.GetValidators()

	require.Equal(t, nil, err2)
	require.Equal(t, validators, validatorsGet)

}
