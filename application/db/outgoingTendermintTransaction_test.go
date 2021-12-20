/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestDeleteOutgoingTendermintTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: Txhash,
	}
	_ = SetOutgoingTendermintTx(tendermintTransaction)
	err = DeleteOutgoingTendermintTx(Txhash)
	require.Nil(t, err)

	db.Close()
}

func TestCountTotalOutgoingTendermintTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"

	expectedTotal, err := CountTotalOutgoingTendermintTx()
	expectedTotal = expectedTotal + 1

	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: Txhash,
	}

	err = SetOutgoingTendermintTx(tendermintTransaction)
	require.Nil(t, err)

	total, err := CountTotalOutgoingTendermintTx()

	require.Nil(t, err)
	require.Equal(t, expectedTotal, total)

	db.Close()
}

func TestIterateOutgoingTmTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	function := func(key []byte, value []byte) error {
		var transactions []OutgoingTendermintTransaction
		var tmTx OutgoingTendermintTransaction
		err := json.Unmarshal(value, &tmTx)
		if err != nil {
			return err
		}
		transactions = append(transactions, tmTx)
		return nil
	}

	err = IterateOutgoingTmTx(function)
	require.Nil(t, err)

	db.Close()
}

func TestNewOutgoingTMTransaction(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: Txhash,
	}
	newTendermintTransaction := NewOutgoingTMTransaction(Txhash)

	require.Equal(t, reflect.TypeOf(tendermintTransaction), reflect.TypeOf(newTendermintTransaction))
	require.Equal(t, tendermintTransaction, newTendermintTransaction)

	db.Close()

}

func TestSetOutgoingTendermintTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
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
	Key := outgoingTendermintTxPrefix.GenerateStoreKey([]byte(tendermintTransaction.TxHash))
	expectedKey := tendermintTransaction.Key()
	require.Equal(t, reflect.TypeOf(Key), reflect.TypeOf(expectedKey))
	require.Equal(t, expectedKey, Key)
}

func TestTendermintOutgoingTransactionValidate(t *testing.T) {
	tmTx := OutgoingTendermintTransaction{}
	require.Equal(t, "OutgoingTendermintTransaction: empty tx hash", tmTx.Validate().Error())

	tmTx.TxHash = "QWERTY"
	require.Equal(t, "OutgoingTendermintTransaction: error decoding tx hash string encoding/hex: invalid byte: U+0051 'Q'", tmTx.Validate().Error())

	tmTx.TxHash = "ABCDEF"
	require.Equal(t, "OutgoingTendermintTransaction: invalid tx hash", tmTx.Validate().Error())

	tmTx.TxHash = "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	err := tmTx.Validate()
	require.Nil(t, err)
}

func TestTendermintOutgoingTransactionValue(t *testing.T) {
	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := OutgoingTendermintTransaction{
		TxHash: Txhash,
	}

	Value, _ := json.Marshal(tendermintTransaction)
	newValue, err := tendermintTransaction.Value()
	require.Nil(t, err)

	require.Equal(t, reflect.TypeOf(Value), reflect.TypeOf(newValue))
	require.Equal(t, Value, newValue)
}

func TestTendermintOutgoingTransactionPrefix(t *testing.T) {
	tendermintTransaction := OutgoingTendermintTransaction{}
	require.Equal(t, reflect.TypeOf(tendermintTransaction.prefix()), reflect.TypeOf(outgoingTendermintTxPrefix))
	require.Equal(t, outgoingTendermintTxPrefix, tendermintTransaction.prefix())
}
