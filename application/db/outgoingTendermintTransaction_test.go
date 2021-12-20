/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/constants"
)

func TestDeleteOutgoingTendermintTx(t *testing.T) {
	db, err := OpenDB(constants.TestDBDir)
	require.Nil(t, err)

	txHash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: txHash,
	}

	_ = SetOutgoingTendermintTx(tendermintTransaction)

	err = DeleteOutgoingTendermintTx(txHash)
	require.Nil(t, err)

	db.Close()
}

func TestCountTotalOutgoingTendermintTx(t *testing.T) {
	db, err := OpenDB(constants.TestDBDir)
	require.Nil(t, err)

	txHash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"

	expectedTotal, err := CountTotalOutgoingTendermintTx()
	require.Nil(t, err)
	expectedTotal++

	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: txHash,
	}

	err = SetOutgoingTendermintTx(tendermintTransaction)
	require.Nil(t, err)

	total, err := CountTotalOutgoingTendermintTx()

	require.Nil(t, err)
	require.Equal(t, expectedTotal, total)

	db.Close()
}

func TestIterateOutgoingTmTx(t *testing.T) {
	db, err := OpenDB(constants.TestDBDir)
	require.Nil(t, err)

	function := func(key []byte, value []byte) error {
		var (
			tmTx OutgoingTendermintTransaction
		)

		innerErr := json.Unmarshal(value, &tmTx)
		if innerErr != nil {
			return innerErr
		}

		return nil
	}

	err = IterateOutgoingTmTx(function)
	require.Nil(t, err)

	db.Close()
}

func TestNewOutgoingTMTransaction(t *testing.T) {
	db, err := OpenDB(constants.TestDBDir)
	require.Nil(t, err)

	txHash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: txHash,
	}

	newTendermintTransaction := NewOutgoingTMTransaction(txHash)

	require.Equal(t, reflect.TypeOf(tendermintTransaction), reflect.TypeOf(newTendermintTransaction))
	require.Equal(t, tendermintTransaction, newTendermintTransaction)

	db.Close()
}

func TestSetOutgoingTendermintTx(t *testing.T) {
	db, err := OpenDB(constants.TestDBDir)
	require.Nil(t, err)

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: Txhash,
	}

	err = SetOutgoingTendermintTx(tendermintTransaction)
	require.Nil(t, err)

	db.Close()
}

func TestTendermintOutgoingTransactionKey(t *testing.T) {
	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F",
	}

	key := outgoingTendermintTxPrefix.GenerateStoreKey([]byte(tendermintTransaction.TxHash))

	expectedKey := tendermintTransaction.Key()
	require.Equal(t, reflect.TypeOf(key), reflect.TypeOf(expectedKey))
	require.Equal(t, expectedKey, key)
}

func TestTendermintOutgoingTransactionValidate(t *testing.T) {
	txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: txhash,
	}

	err := tendermintTransaction.Validate()
	require.Nil(t, err)
}

func TestTendermintOutgoingTransactionValue(t *testing.T) {
	txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: txhash,
	}

	value, _ := json.Marshal(tendermintTransaction)
	newValue, err := tendermintTransaction.Value()
	require.Nil(t, err)

	require.Equal(t, reflect.TypeOf(value), reflect.TypeOf(newValue))
	require.Equal(t, value, newValue)
}

func TestTendermintOutgoingTransactionPrefix(t *testing.T) {
	tendermintTransaction := OutgoingTendermintTransaction{}
	require.Equal(t, reflect.TypeOf(tendermintTransaction.prefix()), reflect.TypeOf(outgoingTendermintTxPrefix))
	require.Equal(t, outgoingTendermintTxPrefix, tendermintTransaction.prefix())
}
