package db

import (
	"bytes"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgraph-io/badger/v3"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
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
		logging.Fatal(err)
	}
	return acc, total
}
