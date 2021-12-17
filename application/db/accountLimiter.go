/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgraph-io/badger/v3"
)

type AccountLimiter struct {
	AccountAddress sdk.AccAddress
	Amount         sdk.Int
}

var _ DBI = &AccountLimiter{}

func (a *AccountLimiter) prefix() storeKeyPrefix {
	return accountLimiterPrefix
}

func (a *AccountLimiter) Key() []byte {
	return a.prefix().GenerateStoreKey(a.AccountAddress.Bytes())
}

func (a *AccountLimiter) Value() ([]byte, error) {
	return json.Marshal(*a)
}

func (a *AccountLimiter) Validate() error {
	if a.Amount.LTE(sdk.ZeroInt()) {
		return fmt.Errorf("invalid amount")
	}
	return nil
}

func GetAccountLimiter(address sdk.AccAddress) (AccountLimiter, error) {
	var acc AccountLimiter
	acc.AccountAddress = address
	acc.Amount = sdk.ZeroInt()
	b, err := get(acc.Key())
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return acc, nil
		} else {
			return acc, err
		}
	}
	err = json.Unmarshal(b, &acc)
	return acc, err
}

func SetAccountLimiter(a AccountLimiter) error {
	return set(&a)
}

func GetTotalTokensWrapped() (sdk.Int, error) {
	total := sdk.ZeroInt()
	err := iterateKeyValues(accountLimiterPrefix.GenerateStoreKey([]byte{}), func(key []byte, value []byte) error {
		var acc AccountLimiter
		jsonErr := json.Unmarshal(value, &acc)
		if jsonErr != nil {
			return jsonErr
		}
		total = total.Add(acc.Amount)
		return nil
	})
	return total, err
}
