package db

import (
	"encoding/json"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddToPendingEthereumIncomingTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	ethInTx := EthereumIncomingTx{
		TxHash:          common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		ProducedToKafka: false,
		Sender:          common.HexToAddress("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		MsgBytes:        nil,
		MsgType:         "",
	}
	err = AddToPendingEthereumIncomingTx(ethInTx)
	require.Equal(t, "empty MsgBytes", err.Error())

	ethInTx.MsgType = bankTypes.MsgSend{}.Type()
	ethInTx.MsgBytes = []byte("Msg")

	err = AddToPendingEthereumIncomingTx(ethInTx)
	require.Equal(t, nil, err)

	db.Close()
}

func TestSetEthereumIncomingTxProduced(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	ethInTx := EthereumIncomingTx{
		TxHash:          common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		ProducedToKafka: false,
		Sender:          common.HexToAddress("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		MsgBytes:        []byte("Msg"),
		MsgType:         bankTypes.MsgSend{}.Type(),
	}
	err = AddToPendingEthereumIncomingTx(ethInTx)
	require.Equal(t, nil, err)

	err = SetEthereumIncomingTxProduced(ethInTx)
	require.Nil(t, err)

	tx, err := GetEthereumIncomingTx(ethInTx.TxHash, true)
	require.Nil(t, err)
	ethInTx.ProducedToKafka = true
	require.Equal(t, ethInTx, tx)

	err = db.Close()
	require.Nil(t, err)

	ethInTx = EthereumIncomingTx{}
	err = SetEthereumIncomingTxProduced(ethInTx)
	require.Equal(t, "DB Closed", err.Error())
}

func TestGetEthereumIncomingTx(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	ethInTx := EthereumIncomingTx{
		TxHash:          common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		ProducedToKafka: false,
		Sender:          common.HexToAddress("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		MsgBytes:        []byte("Msg"),
		MsgType:         bankTypes.MsgSend{}.Type(),
	}
	err = AddToPendingEthereumIncomingTx(ethInTx)
	require.Equal(t, nil, err)

	tx, err := GetEthereumIncomingTx(ethInTx.TxHash, false)
	require.Nil(t, err)
	require.Equal(t, ethInTx, tx)

	err = db.Close()
	require.Nil(t, err)

	ethInTx = EthereumIncomingTx{}
	_, err = GetEthereumIncomingTx(ethInTx.TxHash, false)
	require.Equal(t, "DB Closed", err.Error())
}

func TestEthereumIncomingTxPrefix(t *testing.T) {
	ethInTx := EthereumIncomingTx{}
	require.Equal(t, ethInTx.prefix(), ethereumIncomingTxPrefix)
	require.NotEqual(t, nil, ethInTx.prefix())
}

func TestEthereumIncomingTxKey(t *testing.T) {
	ethInTx := EthereumIncomingTx{
		TxHash:          common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		ProducedToKafka: false,
		Sender:          common.HexToAddress("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		MsgBytes:        []byte("Msg"),
		MsgType:         bankTypes.MsgSend{}.Type(),
	}
	require.Equal(t, ethereumIncomingTxPrefix.GenerateStoreKey(append([]byte{byte(0)}, ethInTx.TxHash.Bytes()...)), ethInTx.Key())
	ethInTx.ProducedToKafka = true
	require.Equal(t, ethereumIncomingTxPrefix.GenerateStoreKey(append([]byte{byte(1)}, ethInTx.TxHash.Bytes()...)), ethInTx.Key())
}

func TestEthereumIncomingTxValue(t *testing.T) {
	ethInTx := EthereumIncomingTx{
		TxHash:          common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		ProducedToKafka: false,
		Sender:          common.HexToAddress("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		MsgBytes:        []byte("Msg"),
		MsgType:         bankTypes.MsgSend{}.Type(),
	}
	b, err := json.Marshal(ethInTx)
	require.Nil(t, err)
	actualBytes, err := ethInTx.Value()
	require.Nil(t, err)
	require.Equal(t, b, actualBytes)
}

func TestEthereumIncomingTxValidate(t *testing.T) {
	ethInTx := EthereumIncomingTx{}
	require.Equal(t, "tx hash is empty", ethInTx.Validate().Error())
	ethInTx.TxHash = common.HexToHash("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2")
	require.Equal(t, "empty MsgBytes", ethInTx.Validate().Error())
	ethInTx.MsgBytes = []byte("Msg")
	require.Equal(t, "invalid msg type", ethInTx.Validate().Error())
	ethInTx.MsgType = bankTypes.MsgSend{}.Type()
	require.Nil(t, ethInTx.Validate())
}

func TestGetProduceToKafkaEthereumIncomingTxs(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	ethInTx := EthereumIncomingTx{
		TxHash:          common.HexToHash("0xbf63e1254d218547057122a7e205233158e5bc48f221bf3a9418b5f7fd576e5f"),
		ProducedToKafka: false,
		Sender:          common.HexToAddress("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		MsgBytes:        []byte("Msg"),
		MsgType:         bankTypes.MsgSend{}.Type(),
	}
	err = AddToPendingEthereumIncomingTx(ethInTx)
	require.Equal(t, nil, err)

	_, err = GetProduceToKafkaEthereumIncomingTxs()
	require.Equal(t, nil, err)

	err = db.Close()
	require.Nil(t, err)

	ethInTx = EthereumIncomingTx{}
	err = SetEthereumIncomingTxProduced(ethInTx)
	require.Equal(t, "DB Closed", err.Error())
}

func TestCheckEthereumIncomingTxProduced(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	ethInTx := EthereumIncomingTx{
		TxHash:          common.HexToHash("0xea85d5fe8c6ee37acd570f801435fb9838a10792f684e60a9f93760cd0b2a4ef"),
		ProducedToKafka: false,
		Sender:          common.HexToAddress("0x679e1f7821bbbb86123c3200a9d4a7f80faa269673357c28b9d6f302454175b2"),
		MsgBytes:        []byte("Msg"),
		MsgType:         bankTypes.MsgSend{}.Type(),
	}

	check := CheckEthereumIncomingTxProduced(ethInTx.TxHash)
	require.Equal(t, false, check)

	err = AddToPendingEthereumIncomingTx(ethInTx)
	require.Equal(t, nil, err)

	check = CheckEthereumIncomingTxProduced(ethInTx.TxHash)
	require.Equal(t, false, check)

	err = db.Close()
	require.Nil(t, err)
}
