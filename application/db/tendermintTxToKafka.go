/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/binary"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgraph-io/badger/v3"
	tmBytes "github.com/tendermint/tendermint/libs/bytes"
)

type TendermintTxToKafka struct {
	TxHash   tmBytes.HexBytes
	MsgIndex uint
	Denom    string
}

var _ KeyValue = &TendermintTxToKafka{}

func (t *TendermintTxToKafka) prefix() storeKeyPrefix {
	return tendermintTxToKafkaPrefix
}

func (t *TendermintTxToKafka) Key() []byte {
	msgIndexBytes := make([]byte, msgUint16Len)
	binary.LittleEndian.PutUint16(msgIndexBytes, uint16(t.MsgIndex))
	denomBytes := []byte(t.Denom)

	return t.prefix().GenerateStoreKey(append(t.TxHash, append(msgIndexBytes, denomBytes...)...))
}

func (t *TendermintTxToKafka) Value() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TendermintTxToKafka) Validate() error {
	if len(t.TxHash.Bytes()) == 0 {
		return ErrEmptyTransaction
	}

	if t.Denom == "" {
		return ErrEmptyDenom
	}

	if err := sdk.ValidateDenom(t.Denom); err != nil {
		return err
	}

	return nil
}

func AddTendermintTxToKafka(database *badger.DB, t TendermintTxToKafka) error {
	return set(database, &t)
}

func GetAllTendermintTxToKafka(database *badger.DB) ([]TendermintTxToKafka, error) {
	var tmTxToKafkaList []TendermintTxToKafka

	err := iterateKeyValues(database, tendermintTxToKafkaPrefix.GenerateStoreKey([]byte{}), func(key []byte, val []byte) error {
		var t TendermintTxToKafka

		innerErr := json.Unmarshal(val, &t)
		if innerErr != nil {
			return innerErr
		}

		tmTxToKafkaList = append(tmTxToKafkaList, t)

		return nil
	})

	return tmTxToKafkaList, err
}

func DeleteTendermintTxToKafka(database *badger.DB, txHash tmBytes.HexBytes, msgIndex uint, denom string) error {
	tmInTx := TendermintTxToKafka{
		TxHash:   txHash,
		MsgIndex: msgIndex,
		Denom:    denom,
	}

	return deleteKV(database, tmInTx.Key())
}
