package db

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
)

type ETHTransaction struct {
	TxHash   common.Hash
	Messages []outgoingTx.WrapTokenMsg
}

func NewETHTransaction(txHash common.Hash, msgs []outgoingTx.WrapTokenMsg) ETHTransaction {
	return ETHTransaction{TxHash: txHash, Messages: msgs}
}

var _ DBI = &ETHTransaction{}

func (ethTx *ETHTransaction) prefix() storeKeyPrefix {
	return ethTransactionPrefix
}

func (ethTx *ETHTransaction) Key() []byte {
	return ethTx.prefix().GenerateStoreKey(ethTx.TxHash.Bytes())
}

func (ethTx *ETHTransaction) Value() ([]byte, error) {
	return json.Marshal(*ethTx)
}

func (ethTx *ETHTransaction) Validate() error {
	// TODO
	if len(ethTx.Messages) == 0 {
		return fmt.Errorf("number of messages for ethHash %s is 0", ethTx.TxHash)
	}
	return nil
}

func SetEthereumTx(ethTransaction ETHTransaction) error {
	return set(&ethTransaction)
}

func IterateEthTx(operation func(key []byte, value []byte) error) error {
	return iterateKeyValues(ethTransactionPrefix.GenerateStoreKey([]byte{}), operation)
}
