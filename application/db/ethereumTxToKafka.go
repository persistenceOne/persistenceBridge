/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
)

type EthereumTxToKafka struct {
	TxHash common.Hash
}

var _ KeyValue = &EthereumTxToKafka{}

func (t *EthereumTxToKafka) prefix() storeKeyPrefix {
	return ethereumTxToKafkaPrefix
}

func (t *EthereumTxToKafka) Key() []byte {
	return t.prefix().GenerateStoreKey(t.TxHash.Bytes())
}

func (t *EthereumTxToKafka) Value() ([]byte, error) {
	return json.Marshal(t)
}

func (t *EthereumTxToKafka) Validate() error {
	if t.TxHash.String() == EthEmptyAddress {
		return ErrEmptyTransaction
	}

	return nil
}

func GetAllEthereumTxToKafka(database *badger.DB) ([]EthereumTxToKafka, error) {
	var ethTxToKafkaList []EthereumTxToKafka

	err := iterateKeyValues(database, ethereumTxToKafkaPrefix.GenerateStoreKey([]byte{}), func(key []byte, val []byte) error {
		var t EthereumTxToKafka

		err := json.Unmarshal(val, &t)
		if err != nil {
			return err
		}

		ethTxToKafkaList = append(ethTxToKafkaList, t)

		return nil
	})

	return ethTxToKafkaList, err
}

func AddEthereumTxToKafka(database *badger.DB, t *EthereumTxToKafka) error {
	return set(database, t)
}

func DeleteEthereumTxToKafka(database *badger.DB, txHash common.Hash) error {
	ethInTx := EthereumTxToKafka{
		TxHash: txHash,
	}

	return deleteKV(database, ethInTx.Key())
}
