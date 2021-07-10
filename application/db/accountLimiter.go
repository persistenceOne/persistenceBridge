package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgraph-io/badger/v3"
	"log"
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
	b, err := get(acc.Key())
	if err != nil {
		return acc, err
	}
	err = json.Unmarshal(b, &acc)
	return acc, err
}

func SetAccountLimiter(a AccountLimiter) error {
	return set(&a)
}

func GetAccountLimiterAndTotal(address sdk.AccAddress) (AccountLimiter, int) {
	total := 0
	acc := AccountLimiter{
		AccountAddress: address,
		Amount:         sdk.ZeroInt(),
	}
	err := iterateKeys(accountLimiterPrefix.GenerateStoreKey([]byte{}), func(key []byte, item *badger.Item) error {
		total = total + 1
		if acc.Amount.Equal(sdk.ZeroInt()) && bytes.Equal(key, acc.Key()) {
			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &acc)
			})
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	return acc, total
}
