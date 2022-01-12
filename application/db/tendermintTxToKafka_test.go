//go:build units

package db

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/persistenceOne/persistenceBridge/utilities/test"
	"github.com/stretchr/testify/require"
)

func TestAddTendermintTxToKafka(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	tmInTx := TendermintTxToKafka{
		TxHash:   nil,
		MsgIndex: 0,
		Denom:    "stake",
	}

	err = AddTendermintTxToKafka(database, tmInTx)
	require.ErrorIs(t, ErrEmptyTransaction, err)

	tmInTx.TxHash, _ = hex.DecodeString("63FD20AC557E7187FDC757D2D0092ECDFE2C062483D633A859DF36687A1E66DC")

	err = AddTendermintTxToKafka(database, tmInTx)
	require.Nil(t, err)
}

func TestGetAllTendermintTxToKafka(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	h, _ := hex.DecodeString("63FD20AC557E7187FDC757D2D0092ECDFE2C062483D633A859DF36687A1E66DC")

	tmInTx := TendermintTxToKafka{
		TxHash:   h,
		MsgIndex: 0,
		Denom:    "stake",
	}

	err = AddTendermintTxToKafka(database, tmInTx)
	require.Nil(t, err)

	txs, err := GetAllTendermintTxToKafka(database)
	require.Nil(t, err)

	exists := false
	for _, tx := range txs {
		if tx.TxHash.String() == tmInTx.TxHash.String() {
			exists = true
			break
		}
	}

	require.Equal(t, true, exists)
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
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	tmInTx := TendermintTxToKafka{
		TxHash:   nil,
		MsgIndex: 0,
		Denom:    "stake",
	}
	tmInTx.TxHash, _ = hex.DecodeString("63FD20AC557E7187FDC757D2D0092ECDFE2C062483D633A859DF36687A1E66DC")

	err = AddTendermintTxToKafka(database, tmInTx)
	require.Equal(t, nil, err)

	err = DeleteTendermintTxToKafka(database, tmInTx.TxHash, tmInTx.MsgIndex, tmInTx.Denom)
	require.Equal(t, nil, err)
}
