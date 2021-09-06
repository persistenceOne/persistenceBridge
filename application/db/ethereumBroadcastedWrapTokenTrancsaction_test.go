package db

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDeleteEthereumTx(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	require.Nil(t, err)

	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   common.Hash{},
		Messages: []outgoingTx.WrapTokenMsg{},
	}
	_ = SetBroadcastedEthereumTx(ethTransaction)

	err = DeleteBroadcastedEthereumTx(common.Hash{})
	require.Nil(t, err)

	db.Close()
}

func TestEthereumBroadcastedWrapTokenTransactionKey(t *testing.T) {
	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   common.Hash{},
		Messages: []outgoingTx.WrapTokenMsg{},
	}

	expectedKey := ethereumBroadcastedWrapTokenTransactionPrefix.GenerateStoreKey(ethTransaction.TxHash.Bytes())
	key := ethTransaction.Key()
	require.Equal(t, expectedKey, key)
}

func TestEthereumBroadcastedWrapTokenTransactionValidate(t *testing.T) {
	Address := common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	wrapTokenMsg := outgoingTx.WrapTokenMsg{
		Address: Address,
		Amount:  big.NewInt(1),
	}
	tx := []outgoingTx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   common.Hash{},
		Messages: tx,
	}
	err := ethTransaction.Validate()
	require.Nil(t, err)

	ethTransaction = EthereumBroadcastedWrapTokenTransaction{
		TxHash:   common.Hash{},
		Messages: []outgoingTx.WrapTokenMsg{},
	}
	err = ethTransaction.Validate()
	expectedErr := fmt.Sprintf("number of messages for ethHash %s is 0", ethTransaction.TxHash)
	require.Equal(t, expectedErr, err.Error())
}

func TestEthereumBroadcastedWrapTokenTransactionValue(t *testing.T) {
	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   common.Hash{},
		Messages: []outgoingTx.WrapTokenMsg{},
	}
	expectedValue, _ := json.Marshal(ethTransaction)
	value, err := ethTransaction.Value()
	require.Nil(t, err)

	require.Equal(t, expectedValue, value)
}

func TestEthereumBroadcastedWrapTokenTransactionPrefix(t *testing.T) {
	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   common.Hash{},
		Messages: []outgoingTx.WrapTokenMsg{},
	}
	Prefix := ethTransaction.prefix()

	require.Equal(t, ethereumBroadcastedWrapTokenTransactionPrefix, Prefix)
}

func TestIterateEthTx(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	require.Nil(t, err)

	function := func(key []byte, value []byte) error {
		var transactions []EthereumBroadcastedWrapTokenTransaction
		var ethTx EthereumBroadcastedWrapTokenTransaction
		err := json.Unmarshal(value, &ethTx)
		require.Nil(t, err)

		transactions = append(transactions, ethTx)
		return nil
	}

	err = IterateBroadcastedEthTx(function)
	require.Nil(t, err)

	db.Close()
}

func TestNewETHTransaction(t *testing.T) {
	Address := common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	wrapTokenMsg := outgoingTx.WrapTokenMsg{
		Address: Address,
		Amount:  big.NewInt(1),
	}
	tx := []outgoingTx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   common.Hash{},
		Messages: tx,
	}

	EthereumBroadcastedWrapTokenTransaction := NewETHTransaction(common.Hash{}, tx)
	err := EthereumBroadcastedWrapTokenTransaction.Validate()
	require.Nil(t, err)

	require.Equal(t, reflect.TypeOf(ethTransaction), reflect.TypeOf(EthereumBroadcastedWrapTokenTransaction))
	require.Equal(t, ethTransaction, EthereumBroadcastedWrapTokenTransaction)
}

func TestSetEthereumTx(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	require.Nil(t, err)

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
	err = SetBroadcastedEthereumTx(ethTransaction)
	require.Nil(t, err)

	db.Close()
}
