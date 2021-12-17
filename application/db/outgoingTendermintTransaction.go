/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v3"
)

type OutgoingTendermintTransaction struct {
	TxHash string
}

func NewOutgoingTMTransaction(txHash string) OutgoingTendermintTransaction {
	return OutgoingTendermintTransaction{TxHash: txHash}
}

var _ DBI = &OutgoingTendermintTransaction{}

func (tmTx *OutgoingTendermintTransaction) prefix() storeKeyPrefix {
	return outgoingTendermintTxPrefix
}

func (tmTx *OutgoingTendermintTransaction) Key() []byte {
	return tmTx.prefix().GenerateStoreKey([]byte(tmTx.TxHash))
}

func (tmTx *OutgoingTendermintTransaction) Value() ([]byte, error) {
	return json.Marshal(*tmTx)
}

func (tmTx *OutgoingTendermintTransaction) Validate() error {
	return nil
}

func DeleteOutgoingTendermintTx(txHash string) error {
	return deleteKV(outgoingTendermintTxPrefix.GenerateStoreKey([]byte(txHash)))
}

func SetOutgoingTendermintTx(tmTransaction OutgoingTendermintTransaction) error {
	return set(&tmTransaction)
}

func IterateOutgoingTmTx(operation func(key []byte, value []byte) error) error {
	return iterateKeyValues(outgoingTendermintTxPrefix.GenerateStoreKey([]byte{}), operation)
}

func CountTotalOutgoingTendermintTx() (int, error) {
	total := 0

	err := iterateKeys(outgoingTendermintTxPrefix.GenerateStoreKey([]byte{}), func(_ []byte, _ *badger.Item) error {
		total++

		return nil
	})

	return total, err
}
