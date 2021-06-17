package db

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TMTransaction struct {
	TxHash   string
	Messages []sdk.Msg
}

func NewTMTransaction(txHash string, msgs []sdk.Msg) TMTransaction {
	return TMTransaction{TxHash: txHash, Messages: msgs}
}

var _ DBI = &TMTransaction{}

func (tmTx *TMTransaction) prefix() storeKeyPrefix {
	return tmTransactionPrefix
}

func (tmTx *TMTransaction) Key() []byte {
	return tmTx.prefix().GenerateStoreKey([]byte(tmTx.TxHash))
}

func (tmTx *TMTransaction) Value() ([]byte, error) {
	return json.Marshal(*tmTx)
}

func (tmTx *TMTransaction) Validate() error {
	// TODO
	if len(tmTx.Messages) == 0 {
		return fmt.Errorf("number of messages for txHash %s is 0", tmTx.TxHash)
	}
	return nil
}

func SetTendermintTx(tmTransaction TMTransaction) error {
	return set(&tmTransaction)
}

func IterateTmTx(operation func(key []byte, value []byte) error) error {
	return iterateKeyValues(tmTransactionPrefix.GenerateStoreKey([]byte{}), operation)
}
