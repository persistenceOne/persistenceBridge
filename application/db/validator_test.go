package db

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSetValidators(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	validatorName := "Binance"
	valoperAddress := "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	require.Nil(t, err)

	err = SetValidator(Validator{
		Address: validatorAddress,
		Name:    validatorName,
	})
	require.Nil(t, err)

	db.Close()
}

func TestGetValidators(t *testing.T) {

	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

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

	db.Close()
}

func TestValidatorKey(t *testing.T) {
	var validatorName string = "Binance"
	var valoperAddress string = "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

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
	var validatorName string = "Binance"
	var valoperAddress string = "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

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
	var validatorName string = "Binance"
	var valoperAddress string = "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

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
