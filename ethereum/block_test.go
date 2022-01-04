/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package ethereum

import (
	"context"
	"math/big"
	"os"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/outgoingtx"
	"github.com/persistenceOne/persistenceBridge/ethereum/contracts"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestCollectEthTx(t *testing.T) {
	configuration.SetConfig(test.GetCmdWithConfig())

	tmAddress, err := casp.GetTendermintAddress()
	require.Nil(t, err)

	configuration.SetPStakeAddress(tmAddress)

	database, err := db.OpenDB(constants.TestDBDir)
	require.Nil(t, err)

	defer database.Close()

	ethereumClient, err := ethclient.Dial(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
	require.Nil(t, err)

	ctx := context.Background()
	tx, _, _ := ethereumClient.TransactionByHash(ctx, common.HexToHash("1f5834f05a156ac8ef9ee1be17b72c1a73e149686364c8fe9509997885ae3409"))
	contract := contracts.LiquidStaking
	encodingConfig := application.MakeEncodingConfig()

	initClientCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TransactionConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authTypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(constants.DefaultPBridgeHome)

	protoCodec := codec.NewProtoCodec(initClientCtx.InterfaceRegistry)

	txHashSuccess := common.HexToHash("0x8e08d80c37c884467b9b48a77e658711615a5cfde43f95fccfb3b95ee66cd6ea")
	address := common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))

	amt := new(big.Int)
	amt.SetInt64(1000)

	wrapTokenMsg := outgoingtx.WrapTokenMsg{
		Address: address,
		Amount:  amt,
	}

	txd := []outgoingtx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := db.OutgoingEthereumTransaction{
		TxHash:   txHashSuccess,
		Messages: txd,
	}

	err = db.SetOutgoingEthereumTx(ethTransaction)
	require.Nil(t, err)

	err = collectEthTx(ctx, ethereumClient, protoCodec, tx, &contract)
	require.Nil(t, err)
}

func TestHandleBlock(t *testing.T) {
	pStakeConfig := configuration.GetAppConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	tmAddress, err := casp.GetTendermintAddress()
	require.Nil(t, err)

	configuration.SetPStakeAddress(tmAddress)

	encodingConfig := application.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TransactionConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authTypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(constants.DefaultPBridgeHome)

	protoCodec := codec.NewProtoCodec(initClientCtx.InterfaceRegistry)

	kafkaProducer := utils.NewProducer(pStakeConfig.Kafka.Brokers, utils.SaramaConfig())
	defer func(kafkaProducer sarama.SyncProducer) {
		innerErr := kafkaProducer.Close()
		if innerErr != nil {
			logging.Error(innerErr)
		}
	}(kafkaProducer)

	database, err := db.OpenDB(constants.TestDBDir)
	require.Nil(t, err)

	defer database.Close()

	ethStatus, err := db.GetEthereumStatus()
	require.Nil(t, err)

	ethereumClient, err := ethclient.Dial(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
	require.Nil(t, err)

	processHeight := big.NewInt(ethStatus.LastCheckHeight + 1)
	ctx := context.Background()

	block, err := ethereumClient.BlockByNumber(ctx, processHeight)
	require.Nil(t, err)

	err = handleBlock(ctx, ethereumClient, block, kafkaProducer, protoCodec)
	require.Nil(t, err)
}

func TestProduceToKafka(t *testing.T) {
	pStakeConfig := configuration.GetAppConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	tmAddress, err := casp.GetTendermintAddress()
	require.Nil(t, err)

	configuration.SetPStakeAddress(tmAddress)

	database, err := db.OpenDB(constants.TestDBDir)
	require.Nil(t, err)

	defer database.Close()

	kafkaProducer := utils.NewProducer(pStakeConfig.Kafka.Brokers, utils.SaramaConfig())
	defer func(kafkaProducer sarama.SyncProducer) {
		innerErr := kafkaProducer.Close()
		if innerErr != nil {
			logging.Error(innerErr)
		}
	}(kafkaProducer)

	txHashSuccess := common.HexToHash("0x8e08d80c37c884467b9b48a77e658711615a5cfde43f95fccfb3b95ee66cd6ea")
	address := common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	amt := new(big.Int)
	amt.SetInt64(1000)

	wrapTokenMsg := outgoingtx.WrapTokenMsg{
		Address: address,
		Amount:  amt,
	}

	txd := []outgoingtx.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := db.OutgoingEthereumTransaction{
		TxHash:   txHashSuccess,
		Messages: txd,
	}

	err = db.SetOutgoingEthereumTx(ethTransaction)
	require.Nil(t, err)

	produceToKafka(kafkaProducer)
}
