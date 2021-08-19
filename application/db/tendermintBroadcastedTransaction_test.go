package db

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDeleteTendermintTx(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := TendermintBroadcastedTransaction{
		TxHash: Txhash,
	}
	_ = SetTendermintTx(tendermintTransaction)
	err = DeleteTendermintTx(Txhash)
	if err != nil {
		t.Fatalf("Error %v", err.Error())
	}
	db.Close()
}

func TestGetTotalTMBroadcastedTx(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	Txhash2 := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC85072"

	tendermintTransaction := TendermintBroadcastedTransaction{
		TxHash: Txhash,
	}

	SetTendermintTx(tendermintTransaction)

	tendermintTransaction = TendermintBroadcastedTransaction{
		TxHash: Txhash2,
	}

	SetTendermintTx(tendermintTransaction)

	total, err := GetTotalTMBroadcastedTx()

	if err != nil {
		t.Fatalf("Error: %v", err.Error())
	}
	if total != 2 {
		t.Fatalf("Error in iterating. Expected Total transactions %v got %v", 2, total)
	}
	db.Close()
}

func TestIterateTmTx(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

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

	err = IterateTmTx(function)
	if err != nil {
		t.Fatalf("Error %v", err.Error())
	}
	db.Close()
}

func TestNewTMTransaction(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := TendermintBroadcastedTransaction{
		TxHash: Txhash,
	}
	newTendermintTransaction := NewTMTransaction(Txhash)

	if reflect.DeepEqual(tendermintTransaction, newTendermintTransaction) == false {
		t.Fatal("New transaction not set correctly")
	}
	db.Close()

}

func TestSetTendermintTx(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := TendermintBroadcastedTransaction{
		TxHash: Txhash,
	}
	err = SetTendermintTx(tendermintTransaction)
	if err != nil {
		t.Fatalf("Error in setting Tendermint Transaction: %v", err.Error())
	}
	db.Close()
}

func TestTendermintBroadcastedTransaction_Key(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := TendermintBroadcastedTransaction{
		TxHash: Txhash,
	}
	if bytes.Compare(tendermintTransaction.prefix().GenerateStoreKey([]byte(tendermintTransaction.TxHash)), tendermintTransaction.Key()) != 0 {
		t.Fatalf("Error in getting Key. Expected %v, Got %v", tendermintTransaction.Key(), tendermintTransaction.prefix().GenerateStoreKey([]byte(tendermintTransaction.TxHash)))
	}
	db.Close()

}

func TestTendermintBroadcastedTransaction_Validate(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := TendermintBroadcastedTransaction{
		TxHash: Txhash,
	}
	err = tendermintTransaction.Validate()
	if err != nil {
		t.Fatalf("Transaction not validated %v", err.Error())
	}
	db.Close()

}

func TestTendermintBroadcastedTransaction_Value(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := TendermintBroadcastedTransaction{
		TxHash: Txhash,
	}

	value, _ := json.Marshal(tendermintTransaction)
	newValue, err := tendermintTransaction.Value()
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	if bytes.Compare(value, newValue) != 0 {
		t.Fatalf("Error in getting Value. Expected %v, Got %v", newValue, value)
	}
	db.Close()

}

func TestTendermintBroadcastedTransaction_prefix(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := "B45A62933F1AC783989F05E6E7C43F9B8D802C41F66A7ED6FEED103CBDC8507F"
	tendermintTransaction := TendermintBroadcastedTransaction{
		TxHash: Txhash,
	}

	prefix := tendermintTransaction.prefix()

	if tendermintBroadcastedTransactionPrefix != prefix {
		t.Fatalf("Error in getting Key. Expected %v, Got %v", EthereumBroadcastedWrapTokenTransactionPrefix, tendermintTransaction.prefix())
	}
	db.Close()

}
