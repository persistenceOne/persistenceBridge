//go:build integration

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
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestCollectEthTx(t *testing.T) {
	configuration.SetConfig(test.GetCmdWithConfig())

	ctx := context.Background()

	tmAddress, err := casp.GetTendermintAddress(ctx)
	require.Nil(t, err)

	ethAddress, err := casp.GetEthAddress(ctx)
	require.Nil(t, err)

	configuration.SetCASPAddresses(tmAddress, ethAddress)

	_, closeFn, err := test.OpenDB(t, db.OpenDB)
	require.Nil(t, err)

	defer closeFn()

	ethereumClient, err := ethclient.Dial(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
	require.Nil(t, err)

	// TODO Need correct tx hash of stake tx of LiquidStaking contract in Ropsten
	//encodingConfig := application.MakeEncodingConfig()
	//initClientCtx := client.Context{}.
	//	WithJSONMarshaler(encodingConfig.Marshaler).
	//	WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
	//	WithTxConfig(encodingConfig.TransactionConfig).
	//	WithLegacyAmino(encodingConfig.Amino).
	//	WithInput(os.Stdin).
	//	WithAccountRetriever(authTypes.AccountRetriever{}).
	//	WithBroadcastMode(flags.BroadcastBlock).
	//	WithHomeDir(constants.DefaultPBridgeHome)
	//protoCodec := codec.NewProtoCodec(initClientCtx.InterfaceRegistry)

	//contract := contracts.LiquidStaking
	TxhashSuccess := common.HexToHash("0xdb95ee137ac5f900db8fef6bd0f1b7f6901ede1e437e5927117b5f5420c00ce0")
	tx, pending, err := ethereumClient.TransactionByHash(ctx, TxhashSuccess)
	require.Nil(t, err)
	require.Equal(t, false, pending)
	require.Equal(t, tx.Hash(), TxhashSuccess)
	//err = collectEthTx(ethereumClient, &ctx, protoCodec, tx, &contract)
	//require.Equal(t, nil, err)

	//err = collectEthTx(ctx, ethereumClient, database, protoCodec, tx, &contract)
	//require.Nil(t, err)
}

func TestHandleBlock(t *testing.T) {
	configuration.SetConfig(test.GetCmdWithConfig())

	ctx := context.Background()

	tmAddress, err := casp.GetTendermintAddress(ctx)
	require.Nil(t, err)

	ethAddress, err := casp.GetEthAddress(ctx)
	require.Nil(t, err)

	configuration.SetCASPAddresses(tmAddress, ethAddress)

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
		WithHomeDir(constants.DefaultPBridgeHome())

	protoCodec := codec.NewProtoCodec(initClientCtx.InterfaceRegistry)

	kafkaProducer := utils.NewProducer(configuration.GetAppConfig().Kafka.Brokers, utils.SaramaConfig())
	defer func(kafkaProducer sarama.SyncProducer) {
		innerErr := kafkaProducer.Close()
		if innerErr != nil {
			logging.Error(innerErr)
		}
	}(kafkaProducer)

	database, closeFn, err := test.OpenDB(t, db.OpenDB)
	defer closeFn()

	require.Nil(t, err)

	ethStatus, err := db.GetEthereumStatus(database)
	require.Nil(t, err)

	ethereumClient, err := ethclient.Dial(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
	require.Nil(t, err)

	processHeight := big.NewInt(ethStatus.LastCheckHeight + 1)

	block, err := ethereumClient.BlockByNumber(ctx, processHeight)
	require.Nil(t, err)

	err = handleBlock(ctx, ethereumClient, database, block, kafkaProducer, protoCodec)
	require.Nil(t, err)
}

func TestProduceToKafka(t *testing.T) {
	pStakeConfig := configuration.GetAppConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	ctx := context.Background()

	tmAddress, err := casp.GetTendermintAddress(ctx)
	require.Nil(t, err)

	ethAddress, err := casp.GetEthAddress(ctx)
	require.Nil(t, err)

	configuration.SetCASPAddresses(tmAddress, ethAddress)

	database, closeFn, err := test.OpenDB(t, db.OpenDB)
	defer closeFn()

	require.Nil(t, err)

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

	wrapTokenMsg := db.WrapTokenMsg{
		Address:       address,
		StakingAmount: amt,
	}

	txd := []db.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := db.OutgoingEthereumTransaction{
		TxHash:   txHashSuccess,
		Messages: txd,
	}

	err = db.SetOutgoingEthereumTx(database, ethTransaction)
	require.Nil(t, err)

	produceToKafka(kafkaProducer, database)
}
