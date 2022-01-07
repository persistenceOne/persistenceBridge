/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetAccountLimiter(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	address1, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	address2, _ := sdk.AccAddressFromBech32("cosmos17p5lujc4d68w5s4usydy60lnh9wx0rrd9ws7mp")

	acc := AccountLimiter{
		AccountAddress: address1,
		Amount:         sdk.OneInt(),
	}

	err = SetAccountLimiter(acc)
	require.Nil(t, err)

	newAccountLimiter1, err := GetAccountLimiter(address1)
	require.Equal(t, acc, newAccountLimiter1)
	require.Nil(t, err)

	newAccountLimiter2, err := GetAccountLimiter(address2)
	require.Nil(t, err)
	require.Equal(t, newAccountLimiter2.AccountAddress, address2)
	require.Equal(t, true, newAccountLimiter2.Amount.Equal(sdk.ZeroInt()))

	db.Close()

	newAccountLimiter2, err = GetAccountLimiter(address2)
	require.Equal(t, "DB Closed", err.Error())
}

func TestAccountLimiterKey(t *testing.T) {
	Address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	acc := AccountLimiter{
		AccountAddress: Address,
		Amount:         sdk.OneInt(),
	}

	key := acc.Key()
	expectedKey := accountLimiterPrefix.GenerateStoreKey(acc.AccountAddress.Bytes())
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
	Address, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	acc := AccountLimiter{
		AccountAddress: Address,
		Amount:         sdk.OneInt(),
	}

	expectedValue, err := acc.Value()
	value, _ := json.Marshal(acc)
	require.Nil(t, err)
	require.Equal(t, expectedValue, value)
}

func TestAccountLimiterPrefix(t *testing.T) {
	acc := AccountLimiter{}

	prefix := acc.prefix()
	require.Equal(t, prefix, accountLimiterPrefix)
	require.NotEqual(t, nil, prefix)
}

func TestSetAccountLimiter(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
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

func TestGetTotalTokensWrapped(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	address1, _ := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	address2, _ := sdk.AccAddressFromBech32("cosmos1ezd6qrpjjj7mgpk8dq2tulnmvzc7mggv7a3ejt")

	acc := AccountLimiter{
		AccountAddress: address1,
		Amount:         sdk.OneInt(),
	}
	err = SetAccountLimiter(acc)
	require.Nil(t, err)

	acc.AccountAddress = address2
	err = SetAccountLimiter(acc)
	require.Nil(t, err)

	total, err := GetTotalTokensWrapped()
	require.Nil(t, err)
	require.Equal(t, sdk.NewInt(2), total)

	db.Close()
}
