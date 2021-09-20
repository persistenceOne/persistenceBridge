package db

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

func TestDeleteEthereumTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash: common.HexToHash("0x134bd3b07e4a39e8e3fa4246533ac7a897ec64c52cbb3a028fe470ce0f1a1375"),
		Messages: []outgoingTx.WrapTokenMsg{{
			Address: common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")),
			Amount:  big.NewInt(1),
		}},
	}
	err = SetBroadcastedEthereumTx(ethTransaction)
	require.Nil(t, err)

	err = DeleteBroadcastedEthereumTx(common.Hash{})
	require.Nil(t, err)

	db.Close()
}

func TestEthereumBroadcastedWrapTokenTransactionKey(t *testing.T) {
	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   common.HexToHash("0x134bd3b07e4a39e8e3fa4246533ac7a897ec64c52cbb3a028fe470ce0f1a1375"),
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
		TxHash:   common.HexToHash("0x134bd3b07e4a39e8e3fa4246533ac7a897ec64c52cbb3a028fe470ce0f1a1375"),
		Messages: tx,
	}
	err := ethTransaction.Validate()
	require.Nil(t, err)

	ethTransaction = EthereumBroadcastedWrapTokenTransaction{
		TxHash:   common.HexToHash("0x134bd3b07e4a39e8e3fa4246533ac7a897ec64c52cbb3a028fe470ce0f1a1375"),
		Messages: []outgoingTx.WrapTokenMsg{},
	}
	err = ethTransaction.Validate()
	require.Equal(t, fmt.Sprintf("number of messages for ethHash %s is 0", ethTransaction.TxHash), err.Error())
	emptyTransaction := EthereumBroadcastedWrapTokenTransaction{}
	require.Equal(t, "tx hash is empty", emptyTransaction.Validate().Error())
}

func TestEthereumBroadcastedWrapTokenTransactionValue(t *testing.T) {
	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash: common.HexToHash("0x134bd3b07e4a39e8e3fa4246533ac7a897ec64c52cbb3a028fe470ce0f1a1375"),
		Messages: []outgoingTx.WrapTokenMsg{{
			Address: common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")),
			Amount:  big.NewInt(1),
		}},
	}
	expectedValue, _ := json.Marshal(ethTransaction)
	value, err := ethTransaction.Value()
	require.Nil(t, err)
	require.Equal(t, expectedValue, value)
}

func TestEthereumBroadcastedWrapTokenTransactionPrefix(t *testing.T) {
	ethTransaction := EthereumBroadcastedWrapTokenTransaction{}
	require.Equal(t, ethereumBroadcastedWrapTokenTransactionPrefix, ethTransaction.prefix())
}

func TestIterateEthTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	function := func(key []byte, value []byte) error {
		var ethTx EthereumBroadcastedWrapTokenTransaction
		err := json.Unmarshal(value, &ethTx)
		require.Nil(t, err)
		return nil
	}

	err = IterateBroadcastedEthTx(function)
	require.Nil(t, err)

	db.Close()
}

func TestNewETHTransaction(t *testing.T) {
	txHash := common.HexToHash("0x134bd3b07e4a39e8e3fa4246533ac7a897ec64c52cbb3a028fe470ce0f1a1375")
	messages := []outgoingTx.WrapTokenMsg{{
		Address: common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")),
		Amount:  big.NewInt(1),
	}}
	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash:   txHash,
		Messages: messages,
	}

	EthereumBroadcastedWrapTokenTransaction := NewETHTransaction(txHash, messages)
	err := EthereumBroadcastedWrapTokenTransaction.Validate()
	require.Nil(t, err)

	require.Equal(t, reflect.TypeOf(ethTransaction), reflect.TypeOf(EthereumBroadcastedWrapTokenTransaction))
	require.Equal(t, ethTransaction, EthereumBroadcastedWrapTokenTransaction)
}

func TestSetEthereumTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	ethTransaction := EthereumBroadcastedWrapTokenTransaction{
		TxHash: common.HexToHash("0x134bd3b07e4a39e8e3fa4246533ac7a897ec64c52cbb3a028fe470ce0f1a1375"),
		Messages: []outgoingTx.WrapTokenMsg{{
			Address: common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")),
			Amount:  big.NewInt(1),
		}},
	}
	err = SetBroadcastedEthereumTx(ethTransaction)
	require.Nil(t, err)

	db.Close()
}
