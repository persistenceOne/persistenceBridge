package rpc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"net/rpc"
)

func AddValidator(validator db.Validator, rpcEndpoint string) ([]db.Validator, error) {
	var result []db.Validator
	client, err := rpc.DialHTTP("tcp", rpcEndpoint)
	defer client.Close()
	if err != nil {
		return result, err
	}

	err = client.Call("ValidatorRPC.AddValidator", validator, &result)
	return result, err
}

func RemoveValidator(validatorAddr sdk.ValAddress, rpcEndpoint string) ([]db.Validator, error) {
	var result []db.Validator
	client, err := rpc.DialHTTP("tcp", rpcEndpoint)
	defer client.Close()
	if err != nil {
		return result, err
	}

	err = client.Call("ValidatorRPC.DeleteValidator", validatorAddr, &result)
	return result, err
}

func ShowValidators(empty string, rpcEndpoint string) ([]db.Validator, error) {
	var result []db.Validator
	client, err := rpc.DialHTTP("tcp", rpcEndpoint)
	defer client.Close()
	if err != nil {
		return result, err
	}

	err = client.Call("ValidatorRPC.GetValidators", empty, &result)
	return result, err

}
