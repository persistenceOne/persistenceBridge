/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type EthereumTxToKafka struct {
	TxHash common.Hash
}

var _ DBI = &EthereumTxToKafka{}

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
	if t.TxHash.String() == "0x0000000000000000000000000000000000000000000000000000000000000000" {
		return fmt.Errorf("tx hash is empty")
	}

	return nil
}

func GetAllEthereumTxToKafka() ([]EthereumTxToKafka, error) {
	var ethTxToKafkaList []EthereumTxToKafka

	err := iterateKeyValues(ethereumTxToKafkaPrefix.GenerateStoreKey([]byte{}), func(key []byte, val []byte) error {
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

func AddEthereumTxToKafka(t *EthereumTxToKafka) error {
	return set(t)
}

func DeleteEthereumTxToKafka(txHash common.Hash) error {
	ethInTx := EthereumTxToKafka{
		TxHash: txHash,
	}

	return deleteKV(ethInTx.Key())
}
