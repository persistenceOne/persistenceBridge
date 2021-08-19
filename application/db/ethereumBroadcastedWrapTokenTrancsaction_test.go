package db

import (
	"bytes"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDeleteEthereumTx(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := common.BytesToHash([]byte("0x800b423ae1dfaf59de9e31fa45ebe0f57268949a8849fc2bd5f054b7c40eb18a"))
	Address := common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	amt := new(big.Int)
	amt.SetInt64(1000)
	wrapTokenMsg := outgoingTx.WrapTokenMsg{
		Address: Address,
		Amount:  amt,
	}
	tx := []outgoingTx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   Txhash,
		Messages: tx,
	}
	_ = SetEthereumTx(ethTransaction)

	err = DeleteEthereumTx(Txhash)
	if err != nil {
		t.Fatalf("Error %v", err.Error())
	}
	db.Close()
}

func TestEthereumBroadcastedWrapTokenTransaction_Key(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := common.BytesToHash([]byte("0x800b423ae1dfaf59de9e31fa45ebe0f57268949a8849fc2bd5f054b7c40eb18a"))
	Address := common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	amt := new(big.Int)
	amt.SetInt64(1000)
	wrapTokenMsg := outgoingTx.WrapTokenMsg{
		Address: Address,
		Amount:  amt,
	}
	tx := []outgoingTx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   Txhash,
		Messages: tx,
	}
	if bytes.Compare(ethTransaction.prefix().GenerateStoreKey(ethTransaction.TxHash.Bytes()), ethTransaction.Key()) != 0 {
		t.Fatalf("Error in getting Key. Expected %v, Got %v", ethTransaction.Key(), ethTransaction.prefix().GenerateStoreKey(ethTransaction.TxHash.Bytes()))
	}
	db.Close()
}

func TestEthereumBroadcastedWrapTokenTransaction_Validate(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := common.BytesToHash([]byte("0x800b423ae1dfaf59de9e31fa45ebe0f57268949a8849fc2bd5f054b7c40eb18a"))
	Address := common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	amt := new(big.Int)
	amt.SetInt64(1000)
	wrapTokenMsg := outgoingTx.WrapTokenMsg{
		Address: Address,
		Amount:  amt,
	}
	tx := []outgoingTx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   Txhash,
		Messages: tx,
	}

	err = ethTransaction.Validate()
	if err != nil {
		t.Fatalf("Error in Validating %v", err.Error())
	}
	db.Close()
}

func TestEthereumBroadcastedWrapTokenTransaction_Value(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := common.BytesToHash([]byte("0x800b423ae1dfaf59de9e31fa45ebe0f57268949a8849fc2bd5f054b7c40eb18a"))
	Address := common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	amt := new(big.Int)
	amt.SetInt64(1000)
	wrapTokenMsg := outgoingTx.WrapTokenMsg{
		Address: Address,
		Amount:  amt,
	}
	tx := []outgoingTx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   Txhash,
		Messages: tx,
	}
	value, _ := json.Marshal(ethTransaction)
	newValue, err := ethTransaction.Value()
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	if bytes.Compare(value, newValue) != 0 {
		t.Fatalf("Error in getting Value. Expected %v, Got %v", newValue, value)
	}
	db.Close()
}

func TestEthereumBroadcastedWrapTokenTransaction_prefix(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := common.BytesToHash([]byte("0x800b423ae1dfaf59de9e31fa45ebe0f57268949a8849fc2bd5f054b7c40eb18a"))
	Address := common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	amt := new(big.Int)
	amt.SetInt64(1000)
	wrapTokenMsg := outgoingTx.WrapTokenMsg{
		Address: Address,
		Amount:  amt,
	}
	tx := []outgoingTx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   Txhash,
		Messages: tx,
	}
	prefix := ethTransaction.prefix()

	if EthereumBroadcastedWrapTokenTransactionPrefix != prefix {
		t.Fatalf("Error in getting Key. Expected %v, Got %v", ethTransaction.Key(), ethTransaction.prefix().GenerateStoreKey(ethTransaction.TxHash.Bytes()))
	}
	db.Close()
}

func TestGetTotalEthBroadcastedTx(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := common.BytesToHash([]byte("0x800b423ae1dfaf59de9e31fa45ebe0f57268949a8849fc2bd5f054b7c40eb420"))
	Txhash2 := common.BytesToHash([]byte("0x800b423ae1dfaf59de9e31fa45ebe0f57268949a8849fc2bd5f054b7c40eb420"))
	Address := common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	amt := new(big.Int)
	amt.SetInt64(1000)
	wrapTokenMsg := outgoingTx.WrapTokenMsg{
		Address: Address,
		Amount:  amt,
	}
	tx := []outgoingTx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   Txhash,
		Messages: tx,
	}

	_ = SetEthereumTx(ethTransaction)

	ethTransaction = EthereumBroadcastedWrapTokenTransaction{
		TxHash:   Txhash2,
		Messages: tx,
	}
	_ = SetEthereumTx(ethTransaction)
	total, err := GetTotalEthBroadcastedTx()
	if err != nil {
		t.Fatalf("Error: %v", err.Error())
	}

	if total != 2 {
		t.Fatalf("Error in iterating. Expected Total transactions %v got %v", 2, total)
	}
	db.Close()
}

func TestIterateEthTx(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	function := func(key []byte, value []byte) error {
		var transactions []EthereumBroadcastedWrapTokenTransaction
		var ethTx EthereumBroadcastedWrapTokenTransaction
		err := json.Unmarshal(value, &ethTx)
		if err != nil {
			return err
		}
		transactions = append(transactions, ethTx)
		return nil
	}

	err = IterateEthTx(function)
	if err != nil {
		t.Fatalf("Error %v", err.Error())
	}
	db.Close()
}

func TestNewETHTransaction(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := common.BytesToHash([]byte("0x800b423ae1dfaf59de9e31fa45ebe0f57268949a8849fc2bd5f054b7c40eb18a"))
	Address := common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	amt := new(big.Int)
	amt.SetInt64(1000)
	wrapTokenMsg := outgoingTx.WrapTokenMsg{
		Address: Address,
		Amount:  amt,
	}
	tx := []outgoingTx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   Txhash,
		Messages: tx,
	}

	EthereumBroadcastedWrapTokenTransaction := NewETHTransaction(Txhash, tx)
	err = EthereumBroadcastedWrapTokenTransaction.Validate()
	if err != nil {
		t.Fatalf("New ethereum tx set not validated: %v", err.Error())
	}
	if reflect.DeepEqual(EthereumBroadcastedWrapTokenTransaction, ethTransaction) == false {
		t.Fatal("New transaction not set correctly")
	}
	db.Close()
}

func TestSetEthereumTx(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	Txhash := common.BytesToHash([]byte("0x800b423ae1dfaf59de9e31fa45ebe0f57268949a8849fc2bd5f054b7c40eb18a"))
	Address := common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	amt := new(big.Int)
	amt.SetInt64(1000)
	wrapTokenMsg := outgoingTx.WrapTokenMsg{
		Address: Address,
		Amount:  amt,
	}
	tx := []outgoingTx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   Txhash,
		Messages: tx,
	}
	err = SetEthereumTx(ethTransaction)
	if err != nil {
		t.Fatalf("Error in setting Transaction: %v", err.Error())
	}
	db.Close()
}
