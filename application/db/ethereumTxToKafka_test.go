//go:build units

package db

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
	"github.com/stretchr/testify/require"
)

func TestAddEthereumTxToKafka(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	ethInTx := &EthereumTxToKafka{
		TxHash: common.Hash{},
	}

	err = AddEthereumTxToKafka(database, ethInTx)
	require.Equal(t, "tx hash is empty", err.Error())

	ethInTx.TxHash = common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2")
	err = AddEthereumTxToKafka(database, ethInTx)
	require.Equal(t, nil, err)
}

func TestGetAllEthereumTxToKafka(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	ethInTx := &EthereumTxToKafka{
		TxHash: common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
	}

	err = AddEthereumTxToKafka(database, ethInTx)
	require.Equal(t, nil, err)

	txs, err := GetAllEthereumTxToKafka(database)
	require.Nil(t, err)

	exists := false
	for _, tx := range txs {
		if tx.TxHash == ethInTx.TxHash {
			exists = true
			break
		}
	}

	require.Equal(t, true, exists)
}

func TestEthereumTxToKafkaPrefix(t *testing.T) {
	ethInTx := EthereumTxToKafka{}
	require.Equal(t, ethInTx.prefix(), ethereumTxToKafkaPrefix)
	require.NotEqual(t, nil, ethInTx.prefix())
}

func TestEthereumTxToKafkaKey(t *testing.T) {
	ethInTx := EthereumTxToKafka{
		TxHash: common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
	}

	require.Equal(t, ethereumTxToKafkaPrefix.GenerateStoreKey(ethInTx.TxHash.Bytes()), ethInTx.Key())
}

func TestEthereumTxToKafkaValue(t *testing.T) {
	ethInTx := EthereumTxToKafka{
		TxHash: common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
	}

	b, err := json.Marshal(ethInTx)
	require.Nil(t, err)

	actualBytes, err := ethInTx.Value()
	require.Nil(t, err)
	require.Equal(t, b, actualBytes)
}

func TestEthereumTxToKafkaValidate(t *testing.T) {
	ethInTx := EthereumTxToKafka{}

	require.Equal(t, "tx hash is empty", ethInTx.Validate().Error())
	ethInTx.TxHash = common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2")
	require.Nil(t, ethInTx.Validate())

}

func TestDeleteEthereumTxToKafka(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	ethInTx := &EthereumTxToKafka{
		TxHash: common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
	}

	err = AddEthereumTxToKafka(database, ethInTx)
	require.Equal(t, nil, err)

	ethInTx.TxHash = common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2")

	err = DeleteEthereumTxToKafka(database, ethInTx.TxHash)
	require.Equal(t, nil, err)
}
