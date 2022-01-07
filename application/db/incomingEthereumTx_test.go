/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddIncomingEthereumTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	ethInTx := IncomingEthereumTx{
		TxHash:   common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		Sender:   common.HexToAddress("0x0000000000000000000000000000000000000001"),
		MsgBytes: nil,
		MsgType:  "",
	}
	err = AddIncomingEthereumTx(ethInTx)
	require.Equal(t, "empty MsgBytes", err.Error())

	ethInTx.MsgType = bankTypes.MsgSend{}.Type()
	ethInTx.MsgBytes = []byte("Msg")

	err = AddIncomingEthereumTx(ethInTx)
	require.Equal(t, nil, err)

	db.Close()
}

func TestGetIncomingEthereumTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	ethInTx := IncomingEthereumTx{
		TxHash:   common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		Sender:   common.HexToAddress("0x0000000000000000000000000000000000000001"),
		MsgBytes: []byte("Msg"),
		MsgType:  bankTypes.MsgSend{}.Type(),
	}
	err = AddIncomingEthereumTx(ethInTx)
	require.Equal(t, nil, err)

	tx, err := GetIncomingEthereumTx(ethInTx.TxHash)
	require.Nil(t, err)
	require.Equal(t, ethInTx, tx)

	err = db.Close()
	require.Nil(t, err)

	ethInTx = IncomingEthereumTx{}
	_, err = GetIncomingEthereumTx(ethInTx.TxHash)
	require.Equal(t, "DB Closed", err.Error())
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

	require.Equal(t, "tx hash is empty", ethInTx.Validate().Error())
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
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)
	require.Equal(t, false, CheckIncomingEthereumTxExists(common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b1")))
	db.Close()
}
