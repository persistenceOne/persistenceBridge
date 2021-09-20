package db

import (
	"encoding/json"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestDeleteTendermintTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := TendermintBroadcastedTransaction{
		TxHash: Txhash,
	}
	_ = SetBroadcastedTendermintTx(tendermintTransaction)
	err = DeleteBroadcastedTendermintTx(Txhash)
	require.Nil(t, err)

	db.Close()
}

func TestGetTotalTMBroadcastedTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"

	expectedTotal, err := GetTotalTMBroadcastedTx()
	expectedTotal = expectedTotal + 1

	tendermintTransaction := TendermintBroadcastedTransaction{
		TxHash: Txhash,
	}

	err = SetBroadcastedTendermintTx(tendermintTransaction)
	require.Nil(t, err)

	total, err := GetTotalTMBroadcastedTx()

	require.Nil(t, err)
	require.Equal(t, expectedTotal, total)

	db.Close()
}

func TestIterateTmTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	function := func(key []byte, value []byte) error {
		var transactions []TendermintBroadcastedTransaction
		var tmTx TendermintBroadcastedTransaction
		err := json.Unmarshal(value, &tmTx)
		if err != nil {
			return err
		}
		transactions = append(transactions, tmTx)
		return nil
	}

	err = IterateBroadcastedTmTx(function)
	require.Nil(t, err)

	db.Close()
}

func TestNewTMTransaction(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := TendermintBroadcastedTransaction{
		TxHash: Txhash,
	}
	newTendermintTransaction := NewTMTransaction(Txhash)

	require.Equal(t, reflect.TypeOf(tendermintTransaction), reflect.TypeOf(newTendermintTransaction))
	require.Equal(t, tendermintTransaction, newTendermintTransaction)

	db.Close()

}

func TestSetTendermintTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := TendermintBroadcastedTransaction{
		TxHash: Txhash,
	}
	err = SetBroadcastedTendermintTx(tendermintTransaction)
	require.Nil(t, err)

	db.Close()
}

func TestTendermintBroadcastedTransactionKey(t *testing.T) {
	tendermintTransaction := TendermintBroadcastedTransaction{
		TxHash: "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F",
	}
	Key := tendermintBroadcastedTransactionPrefix.GenerateStoreKey([]byte(tendermintTransaction.TxHash))
	expectedKey := tendermintTransaction.Key()
	require.Equal(t, reflect.TypeOf(Key), reflect.TypeOf(expectedKey))
	require.Equal(t, expectedKey, Key)
}

func TestTendermintBroadcastedTransactionValidate(t *testing.T) {
	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := TendermintBroadcastedTransaction{
		TxHash: Txhash,
	}
	err := tendermintTransaction.Validate()
	require.Nil(t, err)
}

func TestTendermintBroadcastedTransactionValue(t *testing.T) {
	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := TendermintBroadcastedTransaction{
		TxHash: Txhash,
	}

	Value, _ := json.Marshal(tendermintTransaction)
	newValue, err := tendermintTransaction.Value()
	require.Nil(t, err)

	require.Equal(t, reflect.TypeOf(Value), reflect.TypeOf(newValue))
	require.Equal(t, Value, newValue)
}

func TestTendermintBroadcastedTransactionPrefix(t *testing.T) {
	tendermintTransaction := TendermintBroadcastedTransaction{}
	require.Equal(t, reflect.TypeOf(tendermintTransaction.prefix()), reflect.TypeOf(tendermintBroadcastedTransactionPrefix))
	require.Equal(t, tendermintBroadcastedTransactionPrefix, tendermintTransaction.prefix())
}
