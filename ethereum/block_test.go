package ethereum

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/ethereum/contracts"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"path/filepath"
	"testing"
)

func TestCollectEthTx(t *testing.T) {
	configuration.InitConfig()
	appconfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appconfig)
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)

	configuration.SetPStakeAddress(tmAddress)

	ethereumClient, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/b21966541db246d398fb31402eec2c14")
	require.Equal(t, nil, err)
	ctx := context.Background()
	tx, _, _ := ethereumClient.TransactionByHash(ctx, common.HexToHash("8e08d80c37c884467b9b48a77e658711615a5cfde43f95fccfb3b95ee66cd6ea"))
	contract := contracts.LiquidStaking
	coltx, err := collectEthTx(ethereumClient, &ctx, tx, &contract)
	require.Equal(t, nil, err)
	require.Equal(t, "0x8e08d80c37c884467b9b48a77e658711615a5cfde43f95fccfb3b95ee66cd6ea", coltx.txHash )

}

func TestHandleBlock(t *testing.T){
	pStakeConfig := configuration.InitConfig()
	appconfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appconfig)
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)

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
		err := kafkaProducer.Close()
		if err != nil {
			logging.Error(err)
		}
	}(kafkaProducer)
	dirname, err := os.UserHomeDir()

	database, err := db.OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	defer database.Close()
	ethStatus, err := db.GetEthereumStatus()

	ethereumClient, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/b21966541db246d398fb31402eec2c14")
	processHeight := big.NewInt(ethStatus.LastCheckHeight + 1)
	ctx := context.Background()

	block, err := ethereumClient.BlockByNumber(ctx, processHeight)

	err = handleBlock(ethereumClient, &ctx, block, &kafkaProducer, protoCodec)
	require.Equal(t, nil, err)

}


func TestProduceToKafka(t *testing.T){
	pStakeConfig := configuration.InitConfig()
	appconfig := test.GetCmdWithConfig()
	configuration.SetConfig(&appconfig)
	tmAddress, err := casp.GetTendermintAddress()
	require.Equal(t, nil, err)

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
	var ethTxToTendermintList []ethTxToTendermint
	var ethTxToTM ethTxToTendermint

	msgs := []sdk.Msg{
		&stakingTypes.MsgDelegate{
			DelegatorAddress: "cosmos184u4khydhtlzkfujwpq9dzl34gz8uuh5jnjntz",
			ValidatorAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf",
			Amount:           sdk.Coin{"atom", sdk.NewInt(12)},
		},
		&stakingTypes.MsgUndelegate{
			DelegatorAddress: "cosmos184u4khydhtlzkfujwpq9dzl34gz8uuh5jnjntz",
			ValidatorAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf",
			Amount:           sdk.Coin{"atom", sdk.NewInt(12)},
		},
		&bankTypes.MsgSend{
			FromAddress: "cosmos184u4khydhtlzkfujwpq9dzl34gz8uuh5jnjntz",
			ToAddress:   "cosmos184u4khydhtlzkfujwpq9dzl34gz8uuh5jnjntz",
			Amount:      github_com_cosmos_cosmos_sdk_types.Coins{},
		},

	}

	for _, msg1 := range msgs {
		ethTxToTM.msg = msg1
		ethTxToTendermintList = append(ethTxToTendermintList, ethTxToTM)
	}


	protoCodec := codec.NewProtoCodec(initClientCtx.InterfaceRegistry)

	kafkaProducer := utils.NewProducer(pStakeConfig.Kafka.Brokers, utils.SaramaConfig())
	defer func(kafkaProducer sarama.SyncProducer) {
		err := kafkaProducer.Close()
		if err != nil {
			logging.Error(err)
		}
	}(kafkaProducer)

	produceToKafka(ethTxToTendermintList, &kafkaProducer, protoCodec)
}




