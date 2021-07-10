package rpc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type ValidatorRPC struct{}

func (a *ValidatorRPC) GetValidators(empty string, result *[]db.Validator) error {
	r, err := db.GetValidators()
	*result = r
	return err
}

func (a *ValidatorRPC) GetByValidatorAddress(valAddress sdk.ValAddress, result *db.Validator) error {
	r, err := db.GetValidator(valAddress)
	*result = r
	return err
}

func (a *ValidatorRPC) AddValidator(validator db.Validator, result *[]db.Validator) error {
	err := db.SetValidator(validator)
	if err != nil {
		return err
	}
	r, err := db.GetValidators()
	*result = r
	return err
}

func (a *ValidatorRPC) DeleteValidator(address sdk.ValAddress, result *[]db.Validator) error {
	err := db.DeleteValidator(address)
	if err != nil {
		return err
	}
	r, err := db.GetValidators()
	*result = r
	return err
}

// can add db as an arguement.

func StartServer() {
	validatorRPC := new(ValidatorRPC)
	err := rpc.Register(validatorRPC)
	if err != nil {
		log.Fatal("error registering ValidatorRPC", err)
	}

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":4040")

	if err != nil {
		log.Fatal("Listener error", err)
	}
	log.Printf("serving rpc on port %d", 4040)
	err = http.Serve(listener, nil)

	if err != nil {
		log.Fatal("error serving: ", err)
	}
}
