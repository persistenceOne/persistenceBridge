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

func TestGetAccountLimiter(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	require.Nil(t, err)

	Address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")

	acc := AccountLimiter{
		AccountAddress: Address,
		Amount:         sdk.OneInt(),
	}

	err = SetAccountLimiter(acc)
	require.Nil(t, err)

	opaddress, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	newAccountLimiter, err := GetAccountLimiter(opaddress)
	require.Equal(t, acc, newAccountLimiter)
	require.Nil(t, err)

	db.Close()
}

func TestAccountLimiterKey(t *testing.T) {
	Address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	acc := AccountLimiter{
		AccountAddress: Address,
		Amount:         sdk.OneInt(),
	}

	key := acc.Key()
	expectedKey := accountLimiterPrefix.GenerateStoreKey(acc.AccountAddress.Bytes())
	require.Equal(t, reflect.TypeOf(key), reflect.TypeOf(expectedKey))
	require.Equal(t, expectedKey, key)
	require.NotEqual(t, nil, key)
}

func TestAccountLimiterValidate(t *testing.T) {
	Address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")

	acc1 := AccountLimiter{
		AccountAddress: Address,
		Amount:         sdk.OneInt(),
	}
	acc2 := AccountLimiter{
		AccountAddress: Address,
		Amount:         sdk.ZeroInt(),
	}

	err := acc1.Validate()
	require.Nil(t, err)
	err = acc2.Validate()
	require.Equal(t, "invalid amount", err.Error())
}

func TestAccountLimiterValue(t *testing.T) {
	acc := AccountLimiter{
		AccountAddress: sdk.AccAddress(""),
		Amount:         sdk.OneInt(),
	}

	expectedValue, err := acc.Value()
	value, _ := json.Marshal(acc)
	require.Nil(t, err)
	require.Equal(t, expectedValue, value)
}

func TestAccountLimiterPrefix(t *testing.T) {
	acc := AccountLimiter{
		AccountAddress: sdk.AccAddress(""),
		Amount:         sdk.OneInt(),
	}

	storeKeyPrefix := acc.prefix()
	require.Equal(t, storeKeyPrefix, accountLimiterPrefix)
	require.NotEqual(t, nil, storeKeyPrefix)
}

func TestSetAccountLimiter(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	require.Nil(t, err)

	Address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")

	acc := AccountLimiter{
		AccountAddress: Address,
		Amount:         sdk.OneInt(),
	}

	err = SetAccountLimiter(acc)
	require.Nil(t, err)

	db.Close()
}
