/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package ethereum

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"testing"
)

func TestCollectEthTx(t *testing.T) {
	test.SetTestConfig()
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)
	ethAddress, err := casp.GetEthAddress()
	require.Equal(t, nil, err)
	configuration.SetCASPAddresses(tmAddress, ethAddress)

	database, err := db.OpenDB(constants.TestDbDir)
	require.Nil(t, err)
	defer database.Close()

	ethereumClient, err := ethclient.Dial(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
	require.Equal(t, nil, err)
	ctx := context.Background()
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

}

func TestHandleBlock(t *testing.T) {
	test.SetTestConfig()
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)
	ethAddress, err := casp.GetEthAddress()
	require.Equal(t, nil, err)
	configuration.SetCASPAddresses(tmAddress, ethAddress)

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

	kafkaProducer := utils.NewProducer(configuration.GetAppConfig().Kafka.GetBrokersList(), utils.SaramaConfig())
	defer func(kafkaProducer sarama.SyncProducer) {
		err := kafkaProducer.Close()
		if err != nil {
			logging.Error(err)
		}
	}(kafkaProducer)

	database, err := db.OpenDB(constants.TestDbDir)
	require.Nil(t, err)
	defer database.Close()
	ethStatus, err := db.GetEthereumStatus()

	ethereumClient, err := ethclient.Dial(configuration.GetAppConfig().Ethereum.EthereumEndPoint)
	processHeight := big.NewInt(ethStatus.LastCheckHeight + 1)
	ctx := context.Background()

	block, err := ethereumClient.BlockByNumber(ctx, processHeight)

	err = handleBlock(ethereumClient, &ctx, block, &kafkaProducer, protoCodec)
	require.Equal(t, nil, err)

}

func TestProduceToKafka(t *testing.T) {
	test.SetTestConfig()
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)
	ethAddress, err := casp.GetEthAddress()
	require.Equal(t, nil, err)
	configuration.SetCASPAddresses(tmAddress, ethAddress)

	database, err := db.OpenDB(constants.TestDbDir)
	require.Nil(t, err)
	defer database.Close()
	kafkaProducer := utils.NewProducer(configuration.GetAppConfig().Kafka.GetBrokersList(), utils.SaramaConfig())
	defer func(kafkaProducer sarama.SyncProducer) {
		err := kafkaProducer.Close()
		if err != nil {
			logging.Error(err)
		}
	}(kafkaProducer)

	TxhashSuccess := common.HexToHash("0x8e08d80c37c884467b9b48a77e658711615a5cfde43f95fccfb3b95ee66cd6ea")
	Address := common.BytesToAddress([]byte("0x477573f212a7bdd5f7c12889bd1ad0aa44fb82aa"))
	amt := new(big.Int)
	amt.SetInt64(1000)
	wrapTokenMsg := db.WrapTokenMsg{
		Address:       Address,
		StakingAmount: amt,
	}
	txd := []db.WrapTokenMsg{wrapTokenMsg}

	ethTransaction := db.OutgoingEthereumTransaction{
		TxHash:   TxhashSuccess,
		Messages: txd,
	}

	err = db.SetOutgoingEthereumTx(ethTransaction)
	produceToKafka(&kafkaProducer)
}
