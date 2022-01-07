package db

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddTendermintTxToKafka(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	tmInTx := TendermintTxToKafka{
		TxHash:   nil,
		MsgIndex: 0,
		Denom:    "stake",
	}

	err = AddTendermintTxToKafka(tmInTx)
	require.Equal(t, "empty tx hash", err.Error())

	tmInTx.TxHash, _ = hex.DecodeString("63FD20AC557E7187FDC757D2D0092ECDFE2C062483D633A859DF36687A1E66DC")
	err = AddTendermintTxToKafka(tmInTx)
	require.Equal(t, nil, err)

	db.Close()
}

func TestGetAllTendermintTxToKafka(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	h, _ := hex.DecodeString("63FD20AC557E7187FDC757D2D0092ECDFE2C062483D633A859DF36687A1E66DC")

	tmInTx := TendermintTxToKafka{
		TxHash:   h,
		MsgIndex: 0,
		Denom:    "stake",
	}
	err = AddTendermintTxToKafka(tmInTx)
	require.Equal(t, nil, err)

	txs, err := GetAllTendermintTxToKafka()
	require.Nil(t, err)
	exists := false
	for _, tx := range txs {
		if tx.TxHash.String() == tmInTx.TxHash.String() {
			exists = true
			break
		}
	}
	require.Equal(t, true, exists)

	db.Close()
}

func TestTendermintTxToKafkaPrefix(t *testing.T) {
	tmInTx := TendermintTxToKafka{}
	require.Equal(t, tmInTx.prefix(), tendermintTxToKafkaPrefix)
	require.NotEqual(t, nil, tmInTx.prefix())
}

func TestTendermintTxToKafkaKey(t *testing.T) {
	tmInTx := TendermintTxToKafka{
		TxHash:   nil,
		MsgIndex: 0,
		Denom:    "stake",
	}
	tmInTx.TxHash, _ = hex.DecodeString("63FD20AC557E7187FDC757D2D0092ECDFE2C062483D633A859DF36687A1E66DC")
	msgIndexBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(msgIndexBytes, uint16(0))
	denomBytes := []byte("stake")
	require.Equal(t, tendermintTxToKafkaPrefix.GenerateStoreKey(append(tmInTx.TxHash.Bytes(), append(msgIndexBytes, denomBytes...)...)), tmInTx.Key())
}

func TestTendermintTxToKafkaValue(t *testing.T) {
	tmInTx := TendermintTxToKafka{
		TxHash:   nil,
		MsgIndex: 0,
		Denom:    "stake",
	}
	tmInTx.TxHash, _ = hex.DecodeString("63FD20AC557E7187FDC757D2D0092ECDFE2C062483D633A859DF36687A1E66DC")
	b, err := json.Marshal(tmInTx)
	require.Nil(t, err)
	actualBytes, err := tmInTx.Value()
	require.Nil(t, err)
	require.Equal(t, b, actualBytes)
}

func TestTendermintTxToKafkaValidate(t *testing.T) {
	tmInTx := TendermintTxToKafka{
		TxHash:   nil,
		MsgIndex: 0,
		Denom:    "",
	}

	require.Equal(t, "empty tx hash", tmInTx.Validate().Error())
	tmInTx.TxHash, _ = hex.DecodeString("63FD20AC557E7187FDC757D2D0092ECDFE2C062483D633A859DF36687A1E66DC")
	require.Equal(t, "empty denom", tmInTx.Validate().Error())
	tmInTx.Denom = "stake"
	require.Nil(t, tmInTx.Validate())
}

func TestDeleteTendermintTxToKafka(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	tmInTx := TendermintTxToKafka{
		TxHash:   nil,
		MsgIndex: 0,
		Denom:    "stake",
	}
	tmInTx.TxHash, _ = hex.DecodeString("63FD20AC557E7187FDC757D2D0092ECDFE2C062483D633A859DF36687A1E66DC")

	err = AddTendermintTxToKafka(tmInTx)
	require.Equal(t, nil, err)

	err = DeleteTendermintTxToKafka(tmInTx.TxHash, tmInTx.MsgIndex, tmInTx.Denom)
	require.Equal(t, nil, err)

	db.Close()
}
