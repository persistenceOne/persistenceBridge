package db

import (
	"encoding/json"
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

func DeleteTendermintTx(txHash string) error {
	return deleteKV(tendermintBroadcastedTransactionPrefix.GenerateStoreKey([]byte(txHash)))
}

func SetTendermintTx(tmTransaction TendermintBroadcastedTransaction) error {
	return set(&tmTransaction)
}

func IterateTmTx(operation func(key []byte, value []byte) error) error {
	return iterateKeyValues(tendermintBroadcastedTransactionPrefix.GenerateStoreKey([]byte{}), operation)
}
