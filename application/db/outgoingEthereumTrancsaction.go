/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/persistenceOne/persistenceBridge/application/outgoingtx"
)

var (
	EthEmptyAddress = common.Address{}.String()

	EthEmptyHash       = common.Hash{}
	EthEmptyHashString = common.Hash{}.String()
)

type OutgoingEthereumTransaction struct {
	TxHash   common.Hash
	Messages []outgoingtx.WrapTokenMsg
}

func NewOutgoingETHTransaction(txHash common.Hash, msgs []outgoingtx.WrapTokenMsg) OutgoingEthereumTransaction {
	return OutgoingEthereumTransaction{TxHash: txHash, Messages: msgs}
}

var _ KeyValue = &OutgoingEthereumTransaction{}

func (ethTx *OutgoingEthereumTransaction) prefix() storeKeyPrefix {
	return outgoingEthereumTxPrefix
}

func (ethTx *OutgoingEthereumTransaction) Key() []byte {
	return ethTx.prefix().GenerateStoreKey(ethTx.TxHash.Bytes())
}

func (ethTx *OutgoingEthereumTransaction) Value() ([]byte, error) {
	return json.Marshal(*ethTx)
}

func (ethTx *OutgoingEthereumTransaction) Validate() error {
	if ethTx.TxHash == EthEmptyHash {
		return ErrEmptyTransaction
	}

	if len(ethTx.Messages) == 0 {
		return fmt.Errorf("%w: hash %s", ErrNoTransactionMessages, ethTx.TxHash)
	}

	return nil
}

func DeleteOutgoingEthereumTx(db *badger.DB, txHash common.Hash) error {
	return deleteKV(db, outgoingEthereumTxPrefix.GenerateStoreKey(txHash.Bytes()))
}

func SetOutgoingEthereumTx(db *badger.DB, ethTransaction OutgoingEthereumTransaction) error {
	return set(db, &ethTransaction)
}

func IterateOutgoingEthTx(db *badger.DB, operation func(key []byte, value []byte) error) error {
	return iterateKeyValues(db, outgoingEthereumTxPrefix.GenerateStoreKey([]byte{}), operation)
}
