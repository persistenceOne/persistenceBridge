package db

import (
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgraph-io/badger/v3"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetAccountLimiter(t *testing.T) {
	TestSetAccountLimiter(t)

	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	pStakeConfig := configuration.InitConfig()
	_, err = toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}

	Address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")

	acc := AccountLimiter{
		AccountAddress: Address,
		Amount:         sdk.OneInt(),
	}

	err = SetAccountLimiter(acc)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	opaddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	_, err = GetAccountLimiter(opaddress)
	if err != nil {
		t.Fatalf("Error %s", err.Error())
	}
	db.Close()
}

func TestAccountLimiter_Key(t *testing.T) {
	TestSetAccountLimiter(t)

	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	pStakeConfig := configuration.InitConfig()
	_, err = toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	Address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	AccountLimiter, _ := GetAccountLimiter(Address)
	key := AccountLimiter.Key()

	if bytes.Compare(key, AccountLimiter.prefix().GenerateStoreKey(AccountLimiter.AccountAddress.Bytes())) != 0 {
		t.Fatal("Fetched wrong key")
	}
	db.Close()
}

func TestAccountLimiter_Validate(t *testing.T) {
	TestSetAccountLimiter(t)

	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	pStakeConfig := configuration.InitConfig()
	_, err = toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	Address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	AccountLimiter, _ := GetAccountLimiter(Address)
	err = AccountLimiter.Validate()
	if err != nil {
		t.Fatalf("error: %s", err.Error())
	}
	db.Close()
}

func TestAccountLimiter_Value(t *testing.T) {
	TestSetAccountLimiter(t)

	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	pStakeConfig := configuration.InitConfig()
	_, err = toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	Address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	AccountLimiter, _ := GetAccountLimiter(Address)
	_, err = AccountLimiter.Value()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	db.Close()
}

func TestAccountLimiter_prefix(t *testing.T) {
	TestSetAccountLimiter(t)

	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	pStakeConfig := configuration.InitConfig()
	_, err = toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	Address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	AccountLimiter, _ := GetAccountLimiter(Address)
	storeKeyPrefix := AccountLimiter.prefix()
	if storeKeyPrefix != accountLimiterPrefix {
		t.Fatalf("expected %v got %v", accountLimiterPrefix, storeKeyPrefix)
	}
	db.Close()
}

func TestGetAccountLimiterAndTotal(t *testing.T) {
	TotalAccounts := 0

	TestSetAccountLimiter(t)

	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	pStakeConfig := configuration.InitConfig()
	_, err = toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}

	Address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	newAccountLimiter, total := GetAccountLimiterAndTotal(Address)

	acc := AccountLimiter{
		AccountAddress: Address,
		Amount:         sdk.ZeroInt(),
	}

	_ = SetAccountLimiter(acc)

	err = iterateKeys(accountLimiterPrefix.GenerateStoreKey([]byte{}), func(key []byte, item *badger.Item) error {
		TotalAccounts = TotalAccounts + 1
		if acc.Amount.Equal(sdk.ZeroInt()) && bytes.Equal(key, acc.Key()) {
			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &acc)
			})
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	if !reflect.DeepEqual(newAccountLimiter, acc) {
		t.Fatalf("Expected %v got %v", acc, newAccountLimiter)
	}

	if total != total {
		t.Fatalf("Expected %v got %v", TotalAccounts, total)
	}

	db.Close()
}

func TestSetAccountLimiter(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	pStakeConfig := configuration.InitConfig()
	_, err = toml.DecodeFile(filepath.Join(dirname, "/.persistenceBridge/config.toml"), &pStakeConfig)
	if err != nil {
		log.Fatalf("Error decoding pStakeConfig file: %v\n", err.Error())
	}
	Address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")

	acc := AccountLimiter{
		AccountAddress: Address,
		Amount:         sdk.OneInt(),
	}

	err = SetAccountLimiter(acc)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	db.Close()
}
