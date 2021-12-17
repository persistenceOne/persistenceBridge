/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
)

type OutgoingEthereumTransaction struct {
	TxHash   common.Hash
	Messages []outgoingTx.WrapTokenMsg
}

func NewOutgoingETHTransaction(txHash common.Hash, msgs []outgoingTx.WrapTokenMsg) OutgoingEthereumTransaction {
	return OutgoingEthereumTransaction{TxHash: txHash, Messages: msgs}
}

var _ DBI = &OutgoingEthereumTransaction{}

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
	if ethTx.TxHash.String() == "0x0000000000000000000000000000000000000000000000000000000000000000" {
		return fmt.Errorf("tx hash is empty")
	}
	if len(ethTx.Messages) == 0 {
		return fmt.Errorf("number of messages for ethHash %s is 0", ethTx.TxHash)
	}
	return nil
}

func DeleteOutgoingEthereumTx(txHash common.Hash) error {
	return deleteKV(outgoingEthereumTxPrefix.GenerateStoreKey(txHash.Bytes()))
}

func SetOutgoingEthereumTx(ethTransaction OutgoingEthereumTransaction) error {
	return set(&ethTransaction)
}

func IterateOutgoingEthTx(operation func(key []byte, value []byte) error) error {
	return iterateKeyValues(outgoingEthereumTxPrefix.GenerateStoreKey([]byte{}), operation)
}
