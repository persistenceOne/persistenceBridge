/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

type OutgoingTendermintTransaction struct {
	TxHash string
}

func NewOutgoingTMTransaction(txHash string) OutgoingTendermintTransaction {
	return OutgoingTendermintTransaction{TxHash: txHash}
}

var _ KeyValue = &OutgoingTendermintTransaction{}

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
	if tmTx.TxHash == "" {
		return fmt.Errorf("OutgoingTendermintTransaction: empty tx hash")
	}
	hexBytes, err := hex.DecodeString(tmTx.TxHash)
	if err != nil {
		return fmt.Errorf("OutgoingTendermintTransaction: error decoding tx hash string %v", err)
	}
	if len(hexBytes) != 32 {
		return fmt.Errorf("OutgoingTendermintTransaction: invalid tx hash")
	}
	return nil
}

func DeleteOutgoingTendermintTx(database *badger.DB, txHash string) error {
	return deleteKV(database, outgoingTendermintTxPrefix.GenerateStoreKey([]byte(txHash)))
}

func SetOutgoingTendermintTx(database *badger.DB, tmTransaction OutgoingTendermintTransaction) error {
	return set(database, &tmTransaction)
}

func IterateOutgoingTmTx(database *badger.DB, operation func(key []byte, value []byte) error) error {
	return iterateKeyValues(database, outgoingTendermintTxPrefix.GenerateStoreKey([]byte{}), operation)
}

func CountTotalOutgoingTendermintTx(database *badger.DB) (int, error) {
	total := 0

	err := iterateKeys(database, outgoingTendermintTxPrefix.GenerateStoreKey([]byte{}), func(_ []byte, _ *badger.Item) error {
		total++

		return nil
	})

	return total, err
}
