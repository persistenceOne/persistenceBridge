/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package rpc

import (
	"net/rpc"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/persistenceBridge/application/db"
)

func AddValidator(validator db.Validator, rpcEndpoint string) ([]db.Validator, error) {
	var result []db.Validator

	client, err := rpc.DialHTTP("tcp", rpcEndpoint)
	if err != nil {
		return result, err
	}

	defer client.Close()

	err = client.Call("ValidatorRPC.AddValidator", validator, &result)

	return result, err
}

func RemoveValidator(validatorAddr sdk.ValAddress, rpcEndpoint string) ([]db.Validator, error) {
	var result []db.Validator

	client, err := rpc.DialHTTP("tcp", rpcEndpoint)
	if err != nil {
		return result, err
	}

	defer client.Close()

	err = client.Call("ValidatorRPC.DeleteValidator", validatorAddr, &result)

	return result, err
}

func ShowValidators(empty string, rpcEndpoint string) ([]db.Validator, error) {
	var result []db.Validator

	client, err := rpc.DialHTTP("tcp", rpcEndpoint)
	if err != nil {
		return result, err
	}

	defer client.Close()

	err = client.Call("ValidatorRPC.GetValidators", empty, &result)

	return result, err
}
