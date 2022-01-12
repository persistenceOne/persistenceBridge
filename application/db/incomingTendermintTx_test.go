//go:build units

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgraph-io/badger/v3"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
	"github.com/stretchr/testify/require"
)

func TestAddToPendingIncomingTendermintTx(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	tmInTx := &IncomingTendermintTx{
		TxHash:      []byte("DC6C86075B1466B65BAC2FF08E8A610D"),
		MsgIndex:    0,
		Denom:       "stake",
		FromAddress: "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f",
		Amount:      sdk.NewInt(1),
		Memo:        "",
	}

	err = AddIncomingTendermintTx(database, tmInTx)
	require.Nil(t, err)
}

func TestSetIncomingTendermintTxProduced(t *testing.T) {
	var (
		database *badger.DB
		closeFn  func()
		err      error
	)

	{
		database, closeFn, err = test.OpenDB(t, OpenDB)
		defer closeFn()

		require.Nil(t, err)

		tmInTx := &IncomingTendermintTx{
			TxHash:      []byte("DC6C86075B1466B65BAC2FF08E8A610D"),
			MsgIndex:    0,
			Denom:       "stake",
			FromAddress: "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f",
			Amount:      sdk.NewInt(1),
			Memo:        "",
		}

		err = AddIncomingTendermintTx(database, tmInTx)
		require.Nil(t, err)

		tx, err := GetIncomingTendermintTx(database, tmInTx.TxHash, 0, "stake")
		require.Nil(t, err)
		require.Equal(t, tmInTx, &tx)
	}

	require.Nil(t, database.Close())
}

func TestGetIncomingTendermintTx(t *testing.T) {
	var (
		database *badger.DB
		closeFn  func()
		err      error
	)

	func() {
		database, closeFn, err = test.OpenDB(t, OpenDB)
		defer closeFn()

		require.Nil(t, err)

		tmInTx := &IncomingTendermintTx{
			TxHash:      []byte("DC6C86075B1466B65BAC2FF08E8A610D"),
			MsgIndex:    0,
			Denom:       "stake",
			FromAddress: "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f",
			Amount:      sdk.NewInt(1),
			Memo:        "",
		}

		err = AddIncomingTendermintTx(database, tmInTx)
		require.Nil(t, err)

		tx, err := GetIncomingTendermintTx(database, tmInTx.TxHash, 0, "stake")
		require.Nil(t, err)
		require.Equal(t, tmInTx, &tx)
	}()

	tmInTx := &IncomingTendermintTx{}

	_, err = GetIncomingTendermintTx(database, tmInTx.TxHash, 0, "stake")
	require.ErrorIs(t, err, badger.ErrDBClosed)
}

func TestIncomingTendermintTxPrefix(t *testing.T) {
	tmInTx := IncomingTendermintTx{}
	require.Equal(t, tmInTx.prefix(), incomingTendermintTxPrefix)
	require.NotEqual(t, nil, tmInTx.prefix())
}

func TestIncomingTendermintTxKey(t *testing.T) {
	tmInTx := IncomingTendermintTx{
		TxHash:      []byte("DC6C86075B1466B65BAC2FF08E8A610D"),
		MsgIndex:    0,
		Denom:       "stake",
		FromAddress: "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f",
		Amount:      sdk.NewInt(1),
		Memo:        "",
	}

	msgIndexBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(msgIndexBytes, uint16(0))
	require.Equal(t, incomingTendermintTxPrefix.GenerateStoreKey(append(tmInTx.TxHash, append(msgIndexBytes, []byte(tmInTx.Denom)...)...)), tmInTx.Key())
}

func TestIncomingTendermintTxValue(t *testing.T) {
	tmInTx := IncomingTendermintTx{
		TxHash:      []byte("DC6C86075B1466B65BAC2FF08E8A610D"),
		MsgIndex:    0,
		Denom:       "stake",
		FromAddress: "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f",
		Amount:      sdk.NewInt(1),
		Memo:        "",
	}

	b, err := json.Marshal(tmInTx)
	require.Nil(t, err)

	actualBytes, err := tmInTx.Value()
	require.Nil(t, err)
	require.Equal(t, b, actualBytes)
}

func TestIncomingTendermintTxValidate(t *testing.T) {
	tmInTx := IncomingTendermintTx{}
	require.ErrorIs(t, tmInTx.Validate(), ErrInvalidTransactionHash)

	const (
		txHash     = "DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9"
		denomA     = "a"
		denomStake = "stake"
		from       = "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f"
	)

	h, _ := hex.DecodeString(txHash)
	tmInTx.TxHash = h
	require.Equal(t, "empty denom", tmInTx.Validate().Error())

	tmInTx.Denom = denomA
	require.Equal(t, "invalid denom: a", tmInTx.Validate().Error())

	tmInTx.Denom = denomStake
	require.Equal(t, "from address empty", tmInTx.Validate().Error())

	tmInTx.FromAddress = denomStake
	require.Equal(t, "invalid from address", tmInTx.Validate().Error())

	tmInTx.FromAddress = from
	require.Equal(t, "amount is nil", tmInTx.Validate().Error())

	tmInTx.Amount = sdk.NewInt(-1)
	require.Equal(t, "amount less than 0", tmInTx.Validate().Error())

	tmInTx.Amount = sdk.NewInt(1)
	require.Nil(t, tmInTx.Validate())
}

func TestCheckIncomingTendermintTxExists(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	h, _ := hex.DecodeString("DC6C86075B1466B65BAC2FF08E8A610D")
	require.Equal(t, false, CheckIncomingTendermintTxExists(database, h, 1, "stake"))
}
