package rpc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddValidator(t *testing.T) {
	validatorAddress, err := sdk.ValAddressFromBech32("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf")

	require.Equal(t, nil, err)

	validatorName := "binance"
	rpcEndpoint := "127.0.0.1:4040"
	go StartServer(rpcEndpoint)
	database, err := db.OpenDB("$HOME/persistence/persistenceBridge/application" + "/db")
	validators, err2 := AddValidator(db.Validator{
		Address: validatorAddress,
		Name:   validatorName ,
	}, rpcEndpoint)

	require.Equal(t, nil, err2)

	defer database.Close()
	validatorsGet, err2 := db.GetValidators()

	require.Equal(t, nil, err2)
	require.Equal(t, validators, validatorsGet)
}

func TestRemoveValidator(t *testing.T) {
	validatorAddress, err := sdk.ValAddressFromBech32("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf")

	require.Equal(t, nil, err)

	rpcEndpoint := "127.0.0.1:4040"
	database, err := db.OpenDB("$HOME/persistence/persistenceBridge/application" + "/db")

	require.Equal(t, nil, err)

	validators, err2 := RemoveValidator(validatorAddress, rpcEndpoint)
	require.Equal(t, nil, err2)

	defer database.Close()
	validatorsGet, err2 := db.GetValidators()

	require.Equal(t, nil, err2)
	require.Equal(t, validators, validatorsGet)

}

func TestShowValidators(t *testing.T) {
	rpcEndpoint := "127.0.0.1:4040"

	database, err := db.OpenDB("$HOME/persistence/persistenceBridge/application" + "/db")
	require.Equal(t, nil, err)

	validators, err2 := ShowValidators("", rpcEndpoint)
	require.Equal(t, nil, err2)

	defer database.Close()
	validatorsGet, err2 := db.GetValidators()

	require.Equal(t, nil, err2)
	require.Equal(t, validators, validatorsGet)

}
