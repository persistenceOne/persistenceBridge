package db

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Validator struct {
	Address sdk.ValAddress
	Active  bool
}

var _ DBI = &Validator{}

func (v *Validator) prefix() storeKeyPrefix {
	return validator
}

func (v *Validator) Key() []byte {
	return v.prefix().GenerateStoreKey(v.Address)
}

func (v *Validator) Value() ([]byte, error) {
	return json.Marshal(*v)
}

func (v *Validator) Validate() error {
	// TODO
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
	return Delete(validator.Key())
}
