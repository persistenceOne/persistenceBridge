package db

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestSetValidators(t *testing.T) {
	database, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	validatorName := "Binance"
	valoperAddress := "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	require.Nil(t, err)

	err = SetValidator(Validator{
		Address: validatorAddress,
		Name:    validatorName,
	})
	require.Nil(t, err)

	database.Close()
}

func TestGetValidators(t *testing.T) {
	database, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	validatorName := "Binance"
	valoperAddress := "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	require.Nil(t, err)

	err = SetValidator(Validator{
		Address: validatorAddress,
		Name:    validatorName,
	})
	require.Nil(t, err)

	expectedValidator, err := GetValidator(validatorAddress)
	require.Nil(t, err)

	var testValidator Validator
	testValidator.Address = validatorAddress
	b, err := get(testValidator.Key())
	require.Nil(t, err)

	err = json.Unmarshal(b, &testValidator)

	require.Equal(t, expectedValidator, testValidator)

	validatorSlice, err := GetValidators()

	var testValidators []Validator
	err = iterateKeyValues(validatorPrefix.GenerateStoreKey([]byte{}), func(key []byte, value []byte) error {
		var v Validator
		err := json.Unmarshal(value, &v)
		if err != nil {
			return err
		}
		testValidators = append(testValidators, v)
		return nil
	})

	require.Equal(t, reflect.TypeOf(validatorSlice), reflect.TypeOf(testValidators))
	require.Equal(t, validatorSlice, testValidators)

	database.Close()
}

func TestValidatorKey(t *testing.T) {
	validatorName := "Binance"
	valoperAddress := "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	require.Nil(t, err)

	Validator := Validator{
		Address: validatorAddress,
		Name:    validatorName,
	}

	expectedKey := Validator.Key()
	Key := Validator.prefix().GenerateStoreKey(Validator.Address)

	require.Equal(t, expectedKey, Key)
}

func TestValidatorPrefix(t *testing.T) {
	validatorName := "Binance"
	valoperAddress := "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	require.Nil(t, err)

	Validator := Validator{
		Address: validatorAddress,
		Name:    validatorName,
	}
	Prefix := Validator.prefix()

	require.Equal(t, reflect.TypeOf(Prefix), reflect.TypeOf(validatorPrefix))
	require.Equal(t, Prefix, validatorPrefix)
}

func TestValidatorValue(t *testing.T) {
	validatorName := "Binance"
	valoperAddress := "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	require.Nil(t, err)

	Validator := Validator{
		Address: validatorAddress,
		Name:    validatorName,
	}
	expectedValue, err := Validator.Value()
	require.Nil(t, err)

	value, err := json.Marshal(Validator)
	require.Nil(t, err)

	require.Equal(t, expectedValue, value)
}

func TestDeleteValidator(t *testing.T) {
	database, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	validatorName := "StakingFund"
	valoperAddress := "cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn"

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	require.Nil(t, err)

	err = SetValidator(Validator{
		Address: validatorAddress,
		Name:    validatorName,
	})
	require.Nil(t, err)

	err = DeleteValidator(validatorAddress)
	require.Nil(t, err)

	database.Close()
}

func TestDeleteAllValidators(t *testing.T) {
	database, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	validatorName := "StakingFund"
	valoperAddress := "cosmosvaloper1000ya26q2cmh399q4c5aaacd9lmmdqp90kw2jn"

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	require.Nil(t, err)

	err = SetValidator(Validator{
		Address: validatorAddress,
		Name:    validatorName,
	})
	require.Nil(t, err)

	err = DeleteAllValidators()
	require.Nil(t, err)

	database.Close()
}
