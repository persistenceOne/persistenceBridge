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
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestGetAccountLimiter(t *testing.T) {
	var (
		database *badger.DB
		closeFn  func()
		err      error
	)

	address1, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	address2, _ := sdk.AccAddressFromBech32("cosmos17p5lujc4d68w5s4usydy60lnh9wx0rrd9ws7mp")

	func() {
		database, closeFn, err = test.OpenDB(t, OpenDB)
		defer closeFn()

		require.Nil(t, err)

		acc := AccountLimiter{
			AccountAddress: address1,
			Amount:         sdk.OneInt(),
		}

		err = SetAccountLimiter(database, acc)
		require.Nil(t, err)

		newAccountLimiter1, err := GetAccountLimiter(database, address1)
		require.Equal(t, acc, newAccountLimiter1)
		require.Nil(t, err)

		newAccountLimiter2, err := GetAccountLimiter(database, address2)
		require.Nil(t, err)
		require.Equal(t, newAccountLimiter2.AccountAddress, address2)
		require.Equal(t, true, newAccountLimiter2.Amount.Equal(sdk.ZeroInt()))
	}()

	_, err = GetAccountLimiter(database, address2)

	require.ErrorIs(t, err, badger.ErrDBClosed)
}

func TestAccountLimiterKey(t *testing.T) {
	address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	acc := AccountLimiter{
		AccountAddress: address,
		Amount:         sdk.OneInt(),
	}

	key := acc.Key()
	require.NotEqual(t, nil, key)

	expectedKey := accountLimiterPrefix.GenerateStoreKey(acc.AccountAddress.Bytes())
	require.Equal(t, reflect.TypeOf(key), reflect.TypeOf(expectedKey))
	require.Equal(t, expectedKey, key)
}

func TestAccountLimiterValidate(t *testing.T) {
	address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")

	acc1 := AccountLimiter{
		AccountAddress: address,
		Amount:         sdk.OneInt(),
	}
	acc2 := AccountLimiter{
		AccountAddress: address,
		Amount:         sdk.ZeroInt(),
	}

	err := acc1.Validate()
	require.Nil(t, err)

	err = acc2.Validate()
	require.Equal(t, "invalid amount", err.Error())
}

func TestAccountLimiterValue(t *testing.T) {
	address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	acc := AccountLimiter{
		AccountAddress: address,
		Amount:         sdk.OneInt(),
	}

	expectedValue, err := acc.Value()
	require.Nil(t, err)

	value, _ := json.Marshal(acc)
	require.Equal(t, expectedValue, value)
}

func TestAccountLimiterPrefix(t *testing.T) {
	acc := AccountLimiter{}

	prefix := acc.prefix()
	require.Equal(t, prefix, accountLimiterPrefix)
	require.NotEqual(t, nil, prefix)
}

func TestSetAccountLimiter(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")

	acc := AccountLimiter{
		AccountAddress: address,
		Amount:         sdk.OneInt(),
	}

	err = SetAccountLimiter(database, acc)
	require.Nil(t, err)
}

func TestGetTotalTokensWrapped(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	address1, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	address2, _ := sdk.AccAddressFromBech32("cosmos1ezd6qrpjjj7mgpk8dq2tulnmvzc7mggv7a3ejt")

	acc := AccountLimiter{
		AccountAddress: address1,
		Amount:         sdk.OneInt(),
	}

	err = SetAccountLimiter(database, acc)
	require.Nil(t, err)

	acc.AccountAddress = address2
	err = SetAccountLimiter(database, acc)
	require.Nil(t, err)

	total, err := GetTotalTokensWrapped(database)
	require.Nil(t, err)
	require.Equal(t, sdk.NewInt(2), total)
}
