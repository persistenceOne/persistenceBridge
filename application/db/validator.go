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

var _ KeyValue = &Validator{}

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

func GetValidator(db *badger.DB, address sdk.ValAddress) (Validator, error) {
	var validator Validator

	validator.Address = address

	b, err := get(db, validator.Key())
	if err != nil {
		return Validator{}, err
	}

	err = json.Unmarshal(b, &validator)

	return validator, err
}

func SetValidator(db *badger.DB, validator Validator) error {
	return set(db, &validator)
}

func DeleteValidator(db *badger.DB, address sdk.ValAddress) error {
	var validator Validator
	validator.Address = address

	return deleteKV(db, validator.Key())
}

func GetValidators(db *badger.DB) ([]Validator, error) {
	var validators []Validator

	err := iterateKeyValues(db, validatorPrefix.GenerateStoreKey([]byte{}), func(key []byte, value []byte) error {
		var v Validator

		innerErr := json.Unmarshal(value, &v)
		if innerErr != nil {
			return innerErr
		}

		validators = append(validators, v)

		return nil
	})

	return validators, err
}

func DeleteAllValidators(db *badger.DB) error {
	err := iterateKeys(db, validatorPrefix.GenerateStoreKey([]byte{}), func(key []byte, item *badger.Item) error {
		return deleteKV(db, key)
	})

	return err
}
