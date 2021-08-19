package db

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetAndSetValidators(t *testing.T) {

	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	var validatorName string = "Binance"
	var valoperAddress string = "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	if err != nil {
		t.Errorf("Could not get validators address frkm Bech32 type: %v", err.Error())
	}

	err = SetValidator(Validator{
		Address: validatorAddress,
		Name:    validatorName,
	})
	if err != nil {
		t.Errorf("Error setting validators: %v", err.Error())
	}

	newValidator, err := GetValidator(validatorAddress)
	if err != nil {
		t.Errorf("Error getting validator: %v", err.Error())
	}

	var testValidator Validator
	testValidator.Address = validatorAddress
	b, err := get(testValidator.Key())
	if err != nil {
		t.Errorf("Error: %v", err.Error())
	}
	err = json.Unmarshal(b, &testValidator)

	if reflect.DeepEqual(testValidator, newValidator) == false {
		t.Errorf("Error getting validators: expected %v got %v", testValidator, newValidator)
	}

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

	if reflect.DeepEqual(validatorSlice, testValidators) == false {
		t.Errorf("Error getting validators: expected %v got %v", validatorSlice, testValidators)
	}
	db.Close()
}
func TestValidator_Key(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	var validatorName string = "Binance"
	var valoperAddress string = "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	if err != nil {
		t.Errorf("Could not get validators address frkm Bech32 type: %v", err.Error())
	}

	Validator := Validator{
		Address: validatorAddress,
		Name:    validatorName,
	}
	err = SetValidator(Validator)
	if err != nil {
		t.Errorf("Error setting validators: %v", err.Error())
	}
	Key := Validator.Key()
	if reflect.DeepEqual(Key, Validator.prefix().GenerateStoreKey(Validator.Address)) == false {
		t.Errorf("Error getting validator key: expected %v got %v", Key, Validator.prefix().GenerateStoreKey(Validator.Address))
	}
	db.Close()
}

func TestValidator_prefix(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	var validatorName string = "Binance"
	var valoperAddress string = "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	if err != nil {
		t.Errorf("Could not get validators address frkm Bech32 type: %v", err.Error())
	}

	Validator := Validator{
		Address: validatorAddress,
		Name:    validatorName,
	}
	err = SetValidator(Validator)
	if err != nil {
		t.Errorf("Error setting validators: %v", err.Error())
	}

	if validatorPrefix != Validator.prefix() {
		t.Errorf("Error getting validator prefix: expected %v got %v", validatorPrefix, Validator.prefix())
	}
	db.Close()
}

func TestValidator_Value(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	var validatorName string = "Binance"
	var valoperAddress string = "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf"

	validatorAddress, err := sdk.ValAddressFromBech32(valoperAddress)
	if err != nil {
		t.Errorf("Could not get validators address frkm Bech32 type: %v", err.Error())
	}

	Validator := Validator{
		Address: validatorAddress,
		Name:    validatorName,
	}
	err = SetValidator(Validator)
	if err != nil {
		t.Errorf("Error setting validators: %v", err.Error())
	}
	Value, err := Validator.Value()

	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	newValue, err := json.Marshal(Validator)

	if err != nil {
		t.Fatalf("%v", err.Error())
	}

	if reflect.DeepEqual(Value, newValue) == false {
		t.Errorf("Error getting validator value: expected %v got %v", newValue, Value)
	}

	db.Close()
}
