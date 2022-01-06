/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

func TestDeleteOutgoingEthereumTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	ethTransaction := OutgoingEthereumTransaction{
		TxHash: common.HexToHash("0x134bd3b07e4a39e8e3fa4246533ac7a897ec64c52cbb3a028fe470ce0f1a1375"),
		Messages: []WrapTokenMsg{{
			Address:       common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")),
			StakingAmount: big.NewInt(1),
		}},
	}
	err = SetOutgoingEthereumTx(ethTransaction)
	require.Nil(t, err)

	err = DeleteOutgoingEthereumTx(common.Hash{})
	require.Nil(t, err)

	db.Close()
}

func TestOutgoingEthereumTransactionKey(t *testing.T) {
	ethTransaction := OutgoingEthereumTransaction{
		TxHash:   common.HexToHash("0x134bd3b07e4a39e8e3fa4246533ac7a897ec64c52cbb3a028fe470ce0f1a1375"),
		Messages: []WrapTokenMsg{},
	}

	expectedKey := outgoingEthereumTxPrefix.GenerateStoreKey(ethTransaction.TxHash.Bytes())
	key := ethTransaction.Key()
	require.Equal(t, expectedKey, key)
}

func TestWrapTokenMsgValidate(t *testing.T) {
	wrapTokenMsg := WrapTokenMsg{}

	require.Equal(t, fmt.Errorf("from address empty"), wrapTokenMsg.Validate())

	wrapTokenMsg.FromAddress, _ = sdkTypes.AccAddressFromBech32("cosmos1l44v83h34uv2rz3q4eel8l538v8xfv3uyuvlqs")
	require.Equal(t, fmt.Errorf("invalid eth address"), wrapTokenMsg.Validate())

	wrapTokenMsg.Address = common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	require.Equal(t, fmt.Errorf("invalid tm tx hash"), wrapTokenMsg.Validate())

	wrapTokenMsg.TendermintTxHash, _ = hex.DecodeString("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9")
	require.Equal(t, fmt.Errorf("both amounts nil"), wrapTokenMsg.Validate())

	wrapTokenMsg.StakingAmount = sdkTypes.ZeroInt().BigInt()
	wrapTokenMsg.WrapAmount = sdkTypes.ZeroInt().BigInt()
	require.Equal(t, fmt.Errorf("both amounts zero"), wrapTokenMsg.Validate())

	wrapTokenMsg.WrapAmount = sdkTypes.OneInt().BigInt()
	require.Nil(t, wrapTokenMsg.Validate())
}

func TestOutgoingEthereumTransactionValidate(t *testing.T) {
	wrapTokenMsg := WrapTokenMsg{
		Address: common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")),
	}
	tx := []WrapTokenMsg{wrapTokenMsg}

	ethTransaction := OutgoingEthereumTransaction{
		TxHash:   common.HexToHash("0x134bd3b07e4a39e8e3fa4246533ac7a897ec64c52cbb3a028fe470ce0f1a1375"),
		Messages: tx,
	}
	err := ethTransaction.Validate()
	require.Equal(t, fmt.Sprintf("from address empty"), err.Error())

	wrapTokenMsg.FromAddress, _ = sdkTypes.AccAddressFromBech32("cosmos1l44v83h34uv2rz3q4eel8l538v8xfv3uyuvlqs")
	wrapTokenMsg.TendermintTxHash, _ = hex.DecodeString("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9")
	wrapTokenMsg.StakingAmount = sdkTypes.ZeroInt().BigInt()
	wrapTokenMsg.WrapAmount = sdkTypes.OneInt().BigInt()
	ethTransaction.Messages = []WrapTokenMsg{wrapTokenMsg}
	require.Equal(t, nil, ethTransaction.Validate())

	ethTransaction = OutgoingEthereumTransaction{
		TxHash:   common.HexToHash("0x134bd3b07e4a39e8e3fa4246533ac7a897ec64c52cbb3a028fe470ce0f1a1375"),
		Messages: []WrapTokenMsg{},
	}
	err = ethTransaction.Validate()
	require.Equal(t, fmt.Sprintf("number of messages for ethHash %s is 0", ethTransaction.TxHash), err.Error())
	emptyTransaction := OutgoingEthereumTransaction{}
	require.Equal(t, "tx hash is empty", emptyTransaction.Validate().Error())
}

func TestOutgoingEthereumTransactionValue(t *testing.T) {
	ethTransaction := OutgoingEthereumTransaction{
		TxHash: common.HexToHash("0x134bd3b07e4a39e8e3fa4246533ac7a897ec64c52cbb3a028fe470ce0f1a1375"),
		Messages: []WrapTokenMsg{{
			Address:       common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")),
			StakingAmount: big.NewInt(1),
		}},
	}
	expectedValue, _ := json.Marshal(ethTransaction)
	value, err := ethTransaction.Value()
	require.Nil(t, err)
	require.Equal(t, expectedValue, value)
}

func TestOutgoingEthereumTransactionPrefix(t *testing.T) {
	ethTransaction := OutgoingEthereumTransaction{}
	require.Equal(t, outgoingEthereumTxPrefix, ethTransaction.prefix())
}

func TestIterateEthTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	function := func(key []byte, value []byte) error {
		var ethTx OutgoingEthereumTransaction
		err := json.Unmarshal(value, &ethTx)
		require.Nil(t, err)
		return nil
	}

	err = IterateOutgoingEthTx(function)
	require.Nil(t, err)

	db.Close()
}

func TestNewETHTransaction(t *testing.T) {
	txHash := common.HexToHash("0x134bd3b07e4a39e8e3fa4246533ac7a897ec64c52cbb3a028fe470ce0f1a1375")
	messages := []WrapTokenMsg{{
		Address:       common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")),
		StakingAmount: big.NewInt(1),
	}}
	ethTransaction := OutgoingEthereumTransaction{
		TxHash:   txHash,
		Messages: messages,
	}

	outgoingEthereumTransaction := NewOutgoingETHTransaction(txHash, messages)
	err := outgoingEthereumTransaction.Validate()
	require.Nil(t, err)

	require.Equal(t, reflect.TypeOf(ethTransaction), reflect.TypeOf(outgoingEthereumTransaction))
	require.Equal(t, ethTransaction, outgoingEthereumTransaction)
}

func TestSetEthereumTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	ethTransaction := OutgoingEthereumTransaction{
		TxHash: common.HexToHash("0x134bd3b07e4a39e8e3fa4246533ac7a897ec64c52cbb3a028fe470ce0f1a1375"),
		Messages: []WrapTokenMsg{{
			Address:       common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa")),
			StakingAmount: big.NewInt(1),
		}},
	}
	err = SetOutgoingEthereumTx(ethTransaction)
	require.Nil(t, err)

	db.Close()
}
