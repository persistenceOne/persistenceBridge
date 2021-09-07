package db

import (
	"encoding/binary"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddToPendingTendermintIncomingTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	tmInTx := TendermintIncomingTx{
		ProducedToKafka: false,
		TxHash:          []byte("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9"),
		MsgIndex:        0,
		Denom:           "stake",
		FromAddress:     "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f",
		Amount:          sdk.NewInt(1),
		Memo:            "",
	}
	err = AddToPendingTendermintIncomingTx(tmInTx)
	require.Equal(t, nil, err)

	db.Close()
}

func TestSetTendermintIncomingTxProduced(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	tmInTx := TendermintIncomingTx{
		ProducedToKafka: false,
		TxHash:          []byte("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9"),
		MsgIndex:        0,
		Denom:           "stake",
		FromAddress:     "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f",
		Amount:          sdk.NewInt(1),
		Memo:            "",
	}
	err = AddToPendingTendermintIncomingTx(tmInTx)
	require.Equal(t, nil, err)

	err = SetTendermintIncomingTxProduced(tmInTx)
	require.Nil(t, err)

	tx, err := GetTendermintIncomingTx(tmInTx.TxHash, 0, true, "stake")
	require.Nil(t, err)
	tmInTx.ProducedToKafka = true
	require.Equal(t, tmInTx, tx)

	err = db.Close()
	require.Nil(t, err)

	tmInTx = TendermintIncomingTx{}
	err = SetTendermintIncomingTxProduced(tmInTx)
	require.Equal(t, "DB Closed", err.Error())
}

func TestGetTendermintIncomingTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	tmInTx := TendermintIncomingTx{
		ProducedToKafka: false,
		TxHash:          []byte("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9"),
		MsgIndex:        0,
		Denom:           "stake",
		FromAddress:     "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f",
		Amount:          sdk.NewInt(1),
		Memo:            "",
	}
	err = AddToPendingTendermintIncomingTx(tmInTx)
	require.Equal(t, nil, err)

	tx, err := GetTendermintIncomingTx(tmInTx.TxHash, 0, false, "stake")
	require.Nil(t, err)
	require.Equal(t, tmInTx, tx)

	err = db.Close()
	require.Nil(t, err)

	tmInTx = TendermintIncomingTx{}
	_, err = GetTendermintIncomingTx(tmInTx.TxHash, 0, false, "stake")
	require.Equal(t, "DB Closed", err.Error())
}

func TestTendermintIncomingTxPrefix(t *testing.T) {
	tmInTx := TendermintIncomingTx{}
	require.Equal(t, tmInTx.prefix(), tendermintIncomingTxPrefix)
	require.NotEqual(t, nil, tmInTx.prefix())
}

func TestTendermintIncomingTxKey(t *testing.T) {
	tmInTx := TendermintIncomingTx{
		ProducedToKafka: false,
		TxHash:          []byte("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9"),
		MsgIndex:        0,
		Denom:           "stake",
		FromAddress:     "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f",
		Amount:          sdk.NewInt(1),
		Memo:            "",
	}
	msgIndexBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(msgIndexBytes, uint16(0))
	require.Equal(t, tendermintIncomingTxPrefix.GenerateStoreKey(append([]byte{byte(0)}, append(tmInTx.TxHash, append(msgIndexBytes, []byte(tmInTx.Denom)...)...)...)), tmInTx.Key())
	tmInTx.ProducedToKafka = true
	require.Equal(t, tendermintIncomingTxPrefix.GenerateStoreKey(append([]byte{byte(1)}, append(tmInTx.TxHash, append(msgIndexBytes, []byte(tmInTx.Denom)...)...)...)), tmInTx.Key())
}

func TestTendermintIncomingTxValue(t *testing.T) {
	tmInTx := TendermintIncomingTx{
		ProducedToKafka: false,
		TxHash:          []byte("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9"),
		MsgIndex:        0,
		Denom:           "stake",
		FromAddress:     "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f",
		Amount:          sdk.NewInt(1),
		Memo:            "",
	}
	b, err := json.Marshal(tmInTx)
	require.Nil(t, err)
	actualBytes, err := tmInTx.Value()
	require.Nil(t, err)
	require.Equal(t, b, actualBytes)
}

func TestTendermintIncomingTxValidate(t *testing.T) {
	tmInTx := TendermintIncomingTx{}
	require.Equal(t, "empty tx hash", tmInTx.Validate().Error())
	tmInTx.TxHash = []byte("TendermintIncomingTx")
	require.Equal(t, "invalid tx hash", tmInTx.Validate().Error())
	tmInTx.TxHash = []byte("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9")
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

func TestGetProduceToKafkaTendermintIncomingTxs(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	tmInTx := TendermintIncomingTx{
		ProducedToKafka: false,
		TxHash:          []byte("DC6C86075B1466B65BAC2FF08E8A610DB1C04378695C2D0AD380E997E4277FF9"),
		MsgIndex:        0,
		Denom:           "stake",
		FromAddress:     "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f",
		Amount:          sdk.NewInt(1),
		Memo:            "",
	}
	err = AddToPendingTendermintIncomingTx(tmInTx)
	require.Equal(t, nil, err)

	_, err = GetProduceToKafkaTendermintIncomingTxs()
	require.Equal(t, nil, err)

	err = db.Close()
	require.Nil(t, err)

	tmInTx = TendermintIncomingTx{}
	err = SetTendermintIncomingTxProduced(tmInTx)
	require.Equal(t, "DB Closed", err.Error())
}

func TestCheckTendermintIncomingTxProduced(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	tmInTx := TendermintIncomingTx{
		ProducedToKafka: false,
		TxHash:          []byte("399236E227A688F8584D7FC6800DA0D5C0CAFF0BF48269DD955B634E4CD2CD8F"),
		MsgIndex:        0,
		Denom:           "stake",
		FromAddress:     "cosmos1xa8zh6vjx042rw3kvj9r32sgctm4frpl88rm3f",
		Amount:          sdk.NewInt(1),
		Memo:            "",
	}

	check := CheckTendermintIncomingTxProduced(tmInTx.TxHash, 0, "stake")
	require.Equal(t, false, check)

	err = AddToPendingTendermintIncomingTx(tmInTx)
	require.Equal(t, nil, err)

	check = CheckTendermintIncomingTxProduced(tmInTx.TxHash, 0, "stake")
	require.Equal(t, false, check)

	err = db.Close()
	require.Nil(t, err)
}
