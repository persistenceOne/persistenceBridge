package db

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v3"
)

type TendermintBroadcastedTransaction struct {
	TxHash string
}

func NewTMTransaction(txHash string) TendermintBroadcastedTransaction {
	return TendermintBroadcastedTransaction{TxHash: txHash}
}

var _ DBI = &TendermintBroadcastedTransaction{}

func (tmTx *TendermintBroadcastedTransaction) prefix() storeKeyPrefix {
	return tendermintBroadcastedTransactionPrefix
}

func (tmTx *TendermintBroadcastedTransaction) Key() []byte {
	return tmTx.prefix().GenerateStoreKey([]byte(tmTx.TxHash))
}

func (tmTx *TendermintBroadcastedTransaction) Value() ([]byte, error) {
	return json.Marshal(*tmTx)
}

func (tmTx *TendermintBroadcastedTransaction) Validate() error {
	return nil
}

func DeleteBroadcastedTendermintTx(txHash string) error {
	return deleteKV(tendermintBroadcastedTransactionPrefix.GenerateStoreKey([]byte(txHash)))
}

func SetBroadcastedTendermintTx(tmTransaction TendermintBroadcastedTransaction) error {
	return set(&tmTransaction)
}

func IterateBroadcastedTmTx(operation func(key []byte, value []byte) error) error {
	return iterateKeyValues(tendermintBroadcastedTransactionPrefix.GenerateStoreKey([]byte{}), operation)
}

func GetTotalTMBroadcastedTx() (int, error) {
	total := 0
	err := iterateKeys(tendermintBroadcastedTransactionPrefix.GenerateStoreKey([]byte{}), func(_ []byte, _ *badger.Item) error {
		total = total + 1
		return nil
	})
	return total, err
}
