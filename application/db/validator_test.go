//go:build units

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestSetValidators(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	const (
		validatorName  = "Binance"
		valoperAddress = "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"
	)

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	require.Nil(t, err)

	err = SetValidator(database, Validator{
		Address: validatorAddress,
		Name:    validatorName,
	})
	require.Nil(t, err)
}

func TestGetValidators(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	const (
		validatorName  = "Binance"
		valoperAddress = "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"
	)

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	require.Nil(t, err)

	err = SetValidator(database, Validator{
		Address: validatorAddress,
		Name:    validatorName,
	})
	require.Nil(t, err)

	expectedValidator, err := GetValidator(database, validatorAddress)
	require.Nil(t, err)

	var testValidator Validator
	testValidator.Address = validatorAddress

	b, err := get(database, testValidator.Key())
	require.Nil(t, err)

	err = json.Unmarshal(b, &testValidator)
	require.Nil(t, err)

	require.Equal(t, expectedValidator, testValidator)

	validatorSlice, err := GetValidators(database)
	require.Nil(t, err)

	var testValidators []Validator

	err = iterateKeyValues(database, validatorPrefix.GenerateStoreKey([]byte{}), func(key []byte, value []byte) error {
		var v Validator

		innerErr := json.Unmarshal(value, &v)
		if innerErr != nil {
			return innerErr
		}

		testValidators = append(testValidators, v)

		return nil
	})

	require.Nil(t, err)
	require.Equal(t, reflect.TypeOf(validatorSlice), reflect.TypeOf(testValidators))
	require.Equal(t, validatorSlice, testValidators)
}

func TestValidatorKey(t *testing.T) {
	const (
		validatorName  = "Binance"
		valoperAddress = "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"
	)

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	require.Nil(t, err)

	validator := Validator{
		Address: validatorAddress,
		Name:    validatorName,
	}

	expectedKey := validator.Key()
	key := validator.prefix().GenerateStoreKey(validator.Address)

	require.Equal(t, expectedKey, key)
}

func TestValidatorPrefix(t *testing.T) {
	const (
		validatorName  = "Binance"
		valoperAddress = "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"
	)

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	require.Nil(t, err)

	validator := Validator{
		Address: validatorAddress,
		Name:    validatorName,
	}

	prefix := validator.prefix()

	require.Equal(t, reflect.TypeOf(prefix), reflect.TypeOf(validatorPrefix))
	require.Equal(t, prefix, validatorPrefix)
}

func TestValidatorValue(t *testing.T) {
	validatorName := "Binance"
	valoperAddress := "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	require.Nil(t, err)

	validator := Validator{
		Address: validatorAddress,
		Name:    validatorName,
	}

	expectedValue, err := validator.Value()
	require.Nil(t, err)

	value, err := json.Marshal(validator)
	require.Nil(t, err)

	require.Equal(t, expectedValue, value)
}

func TestDeleteValidator(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	const (
		validatorName  = "StakingFund"
		valoperAddress = "cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn"
	)

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	require.Nil(t, err)

	err = SetValidator(database, Validator{
		Address: validatorAddress,
		Name:    validatorName,
	})
	require.Nil(t, err)

	err = DeleteValidator(database, validatorAddress)
	require.Nil(t, err)
}

func TestDeleteAllValidators(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	validatorName := "StakingFund"
	valoperAddress := "cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn"

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	require.Nil(t, err)

	err = SetValidator(database, Validator{
		Address: validatorAddress,
		Name:    validatorName,
	})
	require.Nil(t, err)

	err = DeleteAllValidators(database)
	require.Nil(t, err)
}
