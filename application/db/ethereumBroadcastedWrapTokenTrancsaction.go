package db

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
)

type EthereumBroadcastedWrapTokenTransaction struct {
	TxHash   common.Hash
	Messages []outgoingTx.WrapTokenMsg
}

func NewETHTransaction(txHash common.Hash, msgs []outgoingTx.WrapTokenMsg) EthereumBroadcastedWrapTokenTransaction {
	return EthereumBroadcastedWrapTokenTransaction{TxHash: txHash, Messages: msgs}
}

var _ DBI = &EthereumBroadcastedWrapTokenTransaction{}

func (ethTx *EthereumBroadcastedWrapTokenTransaction) prefix() storeKeyPrefix {
	return ethereumBroadcastedWrapTokenTransactionPrefix
}

func (ethTx *EthereumBroadcastedWrapTokenTransaction) Key() []byte {
	return ethTx.prefix().GenerateStoreKey(ethTx.TxHash.Bytes())
}

func (ethTx *EthereumBroadcastedWrapTokenTransaction) Value() ([]byte, error) {
	return json.Marshal(*ethTx)
}

func (ethTx *EthereumBroadcastedWrapTokenTransaction) Validate() error {
	if ethTx.TxHash.String() == "0x0000000000000000000000000000000000000000000000000000000000000000" {
		return fmt.Errorf("tx hash is empty")
	}
	if len(ethTx.Messages) == 0 {
		return fmt.Errorf("number of messages for ethHash %s is 0", ethTx.TxHash)
	}
	return nil
}

func DeleteBroadcastedEthereumTx(txHash common.Hash) error {
	return deleteKV(ethereumBroadcastedWrapTokenTransactionPrefix.GenerateStoreKey(txHash.Bytes()))
}

func SetBroadcastedEthereumTx(ethTransaction EthereumBroadcastedWrapTokenTransaction) error {
	return set(&ethTransaction)
}

func IterateBroadcastedEthTx(operation func(key []byte, value []byte) error) error {
	return iterateKeyValues(ethereumBroadcastedWrapTokenTransactionPrefix.GenerateStoreKey([]byte{}), operation)
}
