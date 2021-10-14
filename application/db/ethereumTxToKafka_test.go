package db

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddEthereumTxToKafka(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	ethInTx := EthereumTxToKafka{
		TxHash: common.Hash{},
	}

	err = AddEthereumTxToKafka(ethInTx)
	require.Equal(t, "tx hash is empty", err.Error())

	ethInTx.TxHash = common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2")
	err = AddEthereumTxToKafka(ethInTx)
	require.Equal(t, nil, err)

	db.Close()
}

func TestGetAllEthereumTxToKafka(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	ethInTx := EthereumTxToKafka{
		TxHash: common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
	}
	err = AddEthereumTxToKafka(ethInTx)
	require.Equal(t, nil, err)

	txs, err := GetAllEthereumTxToKafka()
	require.Nil(t, err)
	exists := false
	for _, tx := range txs {
		if tx.TxHash.String() == ethInTx.TxHash.String() {
			exists = true
			break
		}
	}
	require.Equal(t, true, exists)

	db.Close()
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
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	ethInTx := EthereumTxToKafka{
		TxHash: common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
	}

	err = AddEthereumTxToKafka(ethInTx)
	require.Equal(t, nil, err)

	ethInTx.TxHash = common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2")
	err = DeleteEthereumTxToKafka(ethInTx.TxHash)
	require.Equal(t, nil, err)

	db.Close()
}
