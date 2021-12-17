/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgraph-io/badger/v3"
)

type Validator struct {
	Address sdk.ValAddress
	Name    string
}

var _ DBI = &Validator{}

func (v *Validator) prefix() storeKeyPrefix {
	return validatorPrefix
}

func (v *Validator) Key() []byte {
	return v.prefix().GenerateStoreKey(v.Address)
}

func (v *Validator) Value() ([]byte, error) {
	return json.Marshal(*v)
}

func (v *Validator) Validate() error {
	return nil
}

func GetValidator(address sdk.ValAddress) (Validator, error) {
	var validator Validator
	validator.Address = address
	b, err := get(validator.Key())
	if err != nil {
		return Validator{}, err
	}
	err = json.Unmarshal(b, &validator)
	return validator, err
}

func SetValidator(validator Validator) error {
	return set(&validator)
}

func DeleteValidator(address sdk.ValAddress) error {
	var validator Validator
	validator.Address = address
	return deleteKV(validator.Key())
}

func GetValidators() ([]Validator, error) {
	var validators []Validator
	err := iterateKeyValues(validatorPrefix.GenerateStoreKey([]byte{}), func(key []byte, value []byte) error {
		var v Validator
		err := json.Unmarshal(value, &v)
		if err != nil {
			return err
		}
		validators = append(validators, v)
		return nil
	})
	if err != nil {
		return validators, err
	}
	return validators, nil
}

func DeleteAllValidators() error {
	err := iterateKeys(validatorPrefix.GenerateStoreKey([]byte{}), func(key []byte, item *badger.Item) error {
		err := deleteKV(key)
		return err
	})
	return err
}
