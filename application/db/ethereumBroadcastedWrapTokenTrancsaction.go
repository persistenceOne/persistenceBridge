package db

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
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
	// TODO
	if len(ethTx.Messages) == 0 {
		return fmt.Errorf("number of messages for ethHash %s is 0", ethTx.TxHash)
	}
	return nil
}

func DeleteEthereumTx(txHash common.Hash) error {
	return deleteKV(ethereumBroadcastedWrapTokenTransactionPrefix.GenerateStoreKey(txHash.Bytes()))
}

func SetEthereumTx(ethTransaction EthereumBroadcastedWrapTokenTransaction) error {
	return set(&ethTransaction)
}

func IterateEthTx(operation func(key []byte, value []byte) error) error {
	return iterateKeyValues(ethereumBroadcastedWrapTokenTransactionPrefix.GenerateStoreKey([]byte{}), operation)
}

func GetTotalEthBroadcastedTx() (int, error) {
	total := 0
	err := iterateKeys(ethereumBroadcastedWrapTokenTransactionPrefix.GenerateStoreKey([]byte{}), func(_ []byte, _ *badger.Item) error {
		total = total + 1
		return nil
	})
	return total, err
}
