//go:build units

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	"testing"

	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
	"github.com/stretchr/testify/require"
)

func TestAddIncomingEthereumTx(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	ethInTx := &IncomingEthereumTx{
		TxHash:   common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		Sender:   common.HexToAddress("0x0000000000000000000000000000000000000001"),
		MsgBytes: nil,
		MsgType:  "",
	}

	err = AddIncomingEthereumTx(database, ethInTx)
	require.Equal(t, "empty MsgBytes", err.Error())

	ethInTx.MsgType = bankTypes.MsgSend{}.Type()
	ethInTx.MsgBytes = []byte("Msg")

	err = AddIncomingEthereumTx(database, ethInTx)
	require.Nil(t, err)
}

func TestGetIncomingEthereumTx(t *testing.T) {
	var (
		database *badger.DB
		closeFn  func()
		err      error
	)

	func() {
		database, closeFn, err = test.OpenDB(t, OpenDB)
		defer closeFn()

		require.Nil(t, err)

		ethInTx := &IncomingEthereumTx{
			TxHash:   common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
			Sender:   common.HexToAddress("0x0000000000000000000000000000000000000001"),
			MsgBytes: []byte("Msg"),
			MsgType:  bankTypes.MsgSend{}.Type(),
		}

		err = AddIncomingEthereumTx(database, ethInTx)
		require.Nil(t, err)

		tx, err := GetIncomingEthereumTx(database, ethInTx.TxHash)
		require.Nil(t, err)
		require.Equal(t, ethInTx, &tx)
	}()

	ethInTx := &IncomingEthereumTx{}

	_, err = GetIncomingEthereumTx(database, ethInTx.TxHash)
	require.ErrorIs(t, err, badger.ErrDBClosed)
}

func TestIncomingEthereumTxPrefix(t *testing.T) {
	ethInTx := IncomingEthereumTx{}
	require.Equal(t, ethInTx.prefix(), incomingEthereumTxPrefix)
	require.NotEqual(t, nil, ethInTx.prefix())
}

func TestIncomingEthereumTxKey(t *testing.T) {
	ethInTx := IncomingEthereumTx{
		TxHash:   common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		Sender:   common.HexToAddress("0x0000000000000000000000000000000000000001"),
		MsgBytes: []byte("Msg"),
		MsgType:  bankTypes.MsgSend{}.Type(),
	}

	require.Equal(t, incomingEthereumTxPrefix.GenerateStoreKey(ethInTx.TxHash.Bytes()), ethInTx.Key())
}

func TestIncomingEthereumTxValue(t *testing.T) {
	ethInTx := IncomingEthereumTx{
		TxHash:   common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		Sender:   common.HexToAddress("0x0000000000000000000000000000000000000001"),
		MsgBytes: []byte("Msg"),
		MsgType:  bankTypes.MsgSend{}.Type(),
	}

	b, err := json.Marshal(ethInTx)
	require.Nil(t, err)

	actualBytes, err := ethInTx.Value()
	require.Nil(t, err)
	require.Equal(t, b, actualBytes)
}

func TestIncomingEthereumTxValidate(t *testing.T) {
	ethInTx := IncomingEthereumTx{}
	require.ErrorIs(t, ethInTx.Validate(), ErrEmptyTransaction)

	ethInTx.TxHash = common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2")
	require.Equal(t, "empty MsgBytes", ethInTx.Validate().Error())

	ethInTx.MsgBytes = []byte("Msg")
	require.Equal(t, "invalid msg type", ethInTx.Validate().Error())

	ethInTx.MsgType = bankTypes.MsgSend{}.Type()
	require.Equal(t, "invalid sender address", ethInTx.Validate().Error())

	ethInTx.Sender = common.HexToAddress("0x0000000000000000000000000000000000000001")
	require.Nil(t, ethInTx.Validate())
}

func TestCheckIncomingEthereumTxExists(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	require.Equal(t, false, CheckIncomingEthereumTxExists(database, common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b1")))
}
