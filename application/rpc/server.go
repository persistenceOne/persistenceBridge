/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package rpc

import (
	"net/http"
	"net/rpc"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgraph-io/badger/v3"

	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

type ValidatorRPC struct {
	db *badger.DB
}

func NewValidatorRPC(database *badger.DB) *ValidatorRPC {
	return &ValidatorRPC{database}
}

func (a *ValidatorRPC) GetValidators(_ string, result *[]db.Validator) error {
	r, err := db.GetValidators(a.db)

	*result = r

	return err
}

func (a *ValidatorRPC) GetByValidatorAddress(valAddress sdk.ValAddress, result *db.Validator) error {
	r, err := db.GetValidator(a.db, valAddress)

	*result = r

	return err
}

func (a *ValidatorRPC) AddValidator(validator db.Validator, result *[]db.Validator) error {
	err := db.SetValidator(a.db, validator)
	if err != nil {
		return err
	}

	r, err := db.GetValidators(a.db)

	*result = r

	return err
}

func (a *ValidatorRPC) DeleteValidator(address sdk.ValAddress, result *[]db.Validator) error {
	err := db.DeleteValidator(a.db, address)
	if err != nil {
		return err
	}

	r, err := db.GetValidators(a.db)
	*result = r

	return err
}

// fixme: we never stop it, even in tests
// can add db as an argument
func StartServer(rpcEndpoint string, database *badger.DB) {
	validatorRPC := NewValidatorRPC(database)

	err := rpc.Register(validatorRPC)
	if err != nil {
		logging.Fatal("error registering ValidatorRPC:", err)
	}

	rpc.HandleHTTP()

	logging.Info("Starting RPC server on:", rpcEndpoint)

	http.HandleFunc("/status", newStatusHandler(database))
	http.HandleFunc("/validators", newValidatorsHandler(database))

	err = http.ListenAndServe(rpcEndpoint, nil)
	if err != nil {
		logging.Fatal("rpc server:", err)
	}
}
