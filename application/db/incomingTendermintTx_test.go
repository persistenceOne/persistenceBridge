package db

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddToPendingIncomingTendermintTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	txHash, _ := hex.DecodeString("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9")

	tmInTx := IncomingTendermintTx{
		TxHash:      txHash,
		MsgIndex:    0,
		Denom:       "stake",
		FromAddress: "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f",
		Amount:      sdk.NewInt(1),
		Memo:        "",
	}
	err = AddIncomingTendermintTx(tmInTx)
	require.Equal(t, nil, err)

	db.Close()
}

func TestSetIncomingTendermintTxProduced(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	txHash, _ := hex.DecodeString("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9")

	tmInTx := IncomingTendermintTx{
		TxHash:      txHash,
		MsgIndex:    0,
		Denom:       "stake",
		FromAddress: "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f",
		Amount:      sdk.NewInt(1),
		Memo:        "",
	}
	err = AddIncomingTendermintTx(tmInTx)
	require.Equal(t, nil, err)

	tx, err := GetIncomingTendermintTx(tmInTx.TxHash, 0, "stake")
	require.Nil(t, err)
	require.Equal(t, tmInTx, tx)

	err = db.Close()
	require.Nil(t, err)
}

func TestGetIncomingTendermintTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	txHash, _ := hex.DecodeString("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9")

	tmInTx := IncomingTendermintTx{
		TxHash:      txHash,
		MsgIndex:    0,
		Denom:       "stake",
		FromAddress: "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f",
		Amount:      sdk.NewInt(1),
		Memo:        "",
	}
	err = AddIncomingTendermintTx(tmInTx)
	require.Equal(t, nil, err)

	tx, err := GetIncomingTendermintTx(tmInTx.TxHash, 0, "stake")
	require.Nil(t, err)
	require.Equal(t, tmInTx, tx)

	err = db.Close()
	require.Nil(t, err)

	tmInTx = IncomingTendermintTx{}
	_, err = GetIncomingTendermintTx(tmInTx.TxHash, 0, "stake")
	require.Equal(t, "DB Closed", err.Error())
}

func TestIncomingTendermintTxPrefix(t *testing.T) {
	tmInTx := IncomingTendermintTx{}
	require.Equal(t, tmInTx.prefix(), incomingTendermintTxPrefix)
	require.NotEqual(t, nil, tmInTx.prefix())
}

func TestIncomingTendermintTxKey(t *testing.T) {
	txHash, _ := hex.DecodeString("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9")
	tmInTx := IncomingTendermintTx{
		TxHash:      txHash,
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
	txHash, _ := hex.DecodeString("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9")
	tmInTx := IncomingTendermintTx{
		TxHash:      txHash,
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
	require.Equal(t, "IncomingTendermintTx: invalid tx hash", tmInTx.Validate().Error())
	h, _ := hex.DecodeString("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9")
	tmInTx.TxHash = h
	require.Equal(t, "empty denom", tmInTx.Validate().Error())
	tmInTx.Denom = "a"
	require.Equal(t, "invalid denom: a", tmInTx.Validate().Error())
	tmInTx.Denom = "stake"
	require.Equal(t, "from address empty", tmInTx.Validate().Error())
	tmInTx.FromAddress = "stake"
	require.Equal(t, "invalid from address", tmInTx.Validate().Error())
	tmInTx.FromAddress = "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f"
	require.Equal(t, "amount is nil", tmInTx.Validate().Error())
	tmInTx.Amount = sdk.NewInt(-1)
	require.Equal(t, "amount less than 0", tmInTx.Validate().Error())
	tmInTx.Amount = sdk.NewInt(1)
	require.Nil(t, tmInTx.Validate())
}

func TestCheckIncomingTendermintTxExists(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)
	h, _ := hex.DecodeString("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9")
	require.Equal(t, false, CheckIncomingTendermintTxExists(h, 1, "stake"))
	db.Close()
}
