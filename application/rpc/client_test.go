package rpc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/stretchr/testify/require"
	"log"
	"net/rpc"
	"testing"
)

func TestAddValidator(t *testing.T) {
	validatorAddress, err := sdk.ValAddressFromBech32("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf")
	validatorName := "binance"
	rpcEndpoint := "127.0.0.1:4040"
	database, err := db.OpenDB("$HOME/persistence/persistenceBridge/application" + "/db")

	if err != nil {
		log.Printf("Db is already open: %v", err)
		log.Printf("sending rpc")
		validators, err2 := AddValidator(db.Validator{
			Address: validatorAddress,
			Name:   validatorName ,
		}, rpcEndpoint)
		require.Equal(t, nil, err2)

		var result []db.Validator
		client, err := rpc.DialHTTP("tcp", rpcEndpoint)
		defer client.Close()
		require.Equal(t, nil, err)

		err = client.Call("ValidatorRPC.AddValidator", db.Validator{
			Address: validatorAddress,
			Name:   validatorName ,
		}, &result)

		require.Equal(t, validators, result)

	}else {
		defer database.Close()
		err2 := db.SetValidator(db.Validator{
			Address: validatorAddress,
			Name:    validatorName,
		})
		require.Equal(t, nil, err2)

		validators, err2 := db.GetValidators()
		require.Equal(t, nil, err2)

		var result []db.Validator
		client, err := rpc.DialHTTP("tcp", rpcEndpoint)
		defer client.Close()
		require.Equal(t, nil, err)

		err = client.Call("ValidatorRPC.AddValidator", db.Validator{
			Address: validatorAddress,
			Name:   validatorName ,
		}, &result)

		require.Equal(t, validators, result)
	}
}

func TestRemoveValidator(t *testing.T) {
	validatorAddress, err := sdk.ValAddressFromBech32("cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf")
	rpcEndpoint := "127.0.0.1:4040"
	database, err := db.OpenDB("$HOME/persistence/persistenceBridge/application" + "/db")

	if err != nil {
		log.Printf("Db is already open: %v", err)
		log.Printf("sending rpc")
		validators, err2 := RemoveValidator(validatorAddress, rpcEndpoint)
		require.Equal(t, nil, err2)

		var result []db.Validator
		client, err := rpc.DialHTTP("tcp", rpcEndpoint)
		defer client.Close()
		require.Equal(t, nil, err)

		err = client.Call("ValidatorRPC.DeleteValidator", validatorAddress, &result)

		require.Equal(t, validators, result)

	}else {
		defer database.Close()
		err = db.DeleteValidator(validatorAddress)

		require.Equal(t, nil, err)

		validators, err2 := db.GetValidators()
		require.Equal(t, nil, err2)

		var result []db.Validator
		client, err := rpc.DialHTTP("tcp", rpcEndpoint)
		defer client.Close()
		require.Equal(t, nil, err)

		err = client.Call("ValidatorRPC.DeleteValidator", validatorAddress, &result)


		require.Equal(t, validators, result)
	}
}

func TestShowValidators(t *testing.T) {
	rpcEndpoint := "127.0.0.1:4040"
	database, err := db.OpenDB("$HOME/persistence/persistenceBridge/application" + "/db")

	if err != nil {
		log.Printf("Db is already open: %v", err)
		log.Printf("sending rpc")
		validators, err2 := ShowValidators("", rpcEndpoint)
		require.Equal(t, nil, err2)

		var result []db.Validator
		client, err := rpc.DialHTTP("tcp", rpcEndpoint)
		defer client.Close()
		require.Equal(t, nil, err)

		err = client.Call("ValidatorRPC.GetValidators", "", &result)

		require.Equal(t, validators, result)

	}else {
		defer database.Close()

		validators, err2 := db.GetValidators()
		require.Equal(t, nil, err2)

		var result []db.Validator
		client, err := rpc.DialHTTP("tcp", rpcEndpoint)
		defer client.Close()
		require.Equal(t, nil, err)

		err = client.Call("ValidatorRPC.GetValidators", "", &result)


		require.Equal(t, validators, result)
	}
}
