//go:build units

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

	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestDeleteOutgoingTendermintTx(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	const txHash = "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"

	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: txHash,
	}

	err = SetOutgoingTendermintTx(database, tendermintTransaction)
	require.Nil(t, err)

	err = DeleteOutgoingTendermintTx(database, txHash)
	require.Nil(t, err)
}

func TestCountTotalOutgoingTendermintTx(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	const txHash = "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"

	expectedTotal, err := CountTotalOutgoingTendermintTx(database)
	require.Nil(t, err)

	expectedTotal++

	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: txHash,
	}

	err = SetOutgoingTendermintTx(database, tendermintTransaction)
	require.Nil(t, err)

	total, err := CountTotalOutgoingTendermintTx(database)
	require.Nil(t, err)
	require.Equal(t, expectedTotal, total)
}

func TestIterateOutgoingTmTx(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	function := func(key []byte, value []byte) error {
		var tmTx OutgoingTendermintTransaction

		return json.Unmarshal(value, &tmTx)
	}

	err = IterateOutgoingTmTx(database, function)
	require.Nil(t, err)
}

func TestNewOutgoingTMTransaction(t *testing.T) {
	const txHash = "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"

	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: txHash,
	}

	newTendermintTransaction := NewOutgoingTMTransaction(txHash)

	require.Equal(t, reflect.TypeOf(tendermintTransaction), reflect.TypeOf(newTendermintTransaction))
	require.Equal(t, tendermintTransaction, newTendermintTransaction)
}

func TestSetOutgoingTendermintTx(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	const txHash = "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"

	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: txHash,
	}

	err = SetOutgoingTendermintTx(database, tendermintTransaction)
	require.Nil(t, err)
}

func TestTendermintOutgoingTransactionKey(t *testing.T) {
	const txHash = "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"

	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: txHash,
	}

	key := outgoingTendermintTxPrefix.GenerateStoreKey([]byte(tendermintTransaction.TxHash))

	expectedKey := tendermintTransaction.Key()
	require.Equal(t, reflect.TypeOf(key), reflect.TypeOf(expectedKey))
	require.Equal(t, expectedKey, key)
}

func TestTendermintOutgoingTransactionValidate(t *testing.T) {
	const txHash = "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"

	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: txHash,
	}

	err := tendermintTransaction.Validate()
	require.Nil(t, err)
}

func TestTendermintOutgoingTransactionValue(t *testing.T) {
	const txHash = "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"

	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: txHash,
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
