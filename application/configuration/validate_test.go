package configuration

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestValidateEthereum(t *testing.T) {
	ethConfig := ethereumConfig{}
	ethConfig.EthereumEndPoint = "invalid url"
	ethConfig.GasLimit = 0
	ethConfig.GasFeeCap = -1
	err := ethConfig.validate()
	require.Equal(t, fmt.Errorf("invalid EthereumEndPoint: parse \"invalid url\": invalid URI for request"), err)
	ethConfig.EthereumEndPoint = "ws://127.0.0.1:8546"
	err = ethConfig.validate()
	require.Equal(t, fmt.Errorf("invalid eth gas limit"), err)
	ethConfig.GasLimit = 1
	err = ethConfig.validate()
	require.Equal(t, fmt.Errorf("invalid eth gas fee cap"), err)
	ethConfig.GasFeeCap = 1
	err = ethConfig.validate()
	require.Equal(t, fmt.Errorf("empty liquid staking contract address"), err)
	ethConfig.LiquidStakingAddress = "0x0000000000000000000000000000000000000001"
	err = ethConfig.validate()
	require.Equal(t, fmt.Errorf("empty token wrapper contract address"), err)
	ethConfig.TokenWrapperAddress = "0x0000000000000000000000000000000000000001"
	err = ethConfig.validate()
	require.Equal(t, fmt.Errorf("bridgeAdminAddress is empty"), err)
	ethConfig.bridgeAdminAddress = common.HexToAddress("0x0000000000000000000000000000000000000001")
	ethConfig.BalanceCheckPeriod = 1
	err = ethConfig.validate()
	require.Equal(t, fmt.Errorf("invalid ethereum balance alert configuration"), err)
	ethConfig.AlertAmount = 1000000000

	err = ethConfig.validate()
	require.Nil(t, err)
}
func TestValidateTendermint(t *testing.T) {
	tmConfig := tendermintConfig{}
	err := tmConfig.validate()
	require.Equal(t, fmt.Errorf("account prefix cannot be empty"), err)
	tmConfig.AccountPrefix = "cosmos"

	err = tmConfig.validate()
	require.Equal(t, fmt.Errorf("denom cannot be empty"), err)
	tmConfig.Denom = "stake"

	err = tmConfig.validate()
	require.Equal(t, fmt.Errorf("invalied tendermint gas price strconv.ParseFloat: parsing \"\": invalid syntax"), err)
	tmConfig.GasPrice = "0.025"

	err = tmConfig.validate()
	require.Equal(t, fmt.Errorf("tendermint gas adjustment should be greater than 1 (possibly 1.5, current: 0)"), err)
	tmConfig.GasAdjustment = 1.5

	tmConfig.MinimumWrapAmount = -1
	err = tmConfig.validate()
	require.Equal(t, fmt.Errorf("minimum wrap amount cannot be less than 0"), err)
	tmConfig.MinimumWrapAmount = 10

	err = tmConfig.validate()
	require.Equal(t, fmt.Errorf("chain id cannot be empty"), err)
	tmConfig.ChainID = "test"

	err = tmConfig.validate()
	require.Equal(t, fmt.Errorf("invalid tendermint node: parse \"\": empty url"), err)
	tmConfig.Node = "http://127.0.0.1:26657"

	err = tmConfig.validate()
	require.Equal(t, fmt.Errorf("invalid broadcast mode"), err)
	tmConfig.BroadcastMode = "sync"

	err = tmConfig.validate()
	require.Equal(t, "wrapAddress empty", err.Error())
	tmConfig.wrapAddress = "wrong address"

	err = tmConfig.validate()
	require.Equal(t, "decoding bech32 failed: invalid character in string: ' '", err.Error())
	tmConfig.wrapAddress = "cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u"

	err = tmConfig.validate()
	require.Equal(t, "tendermint chain avg block time cannot be 0", err.Error())
	tmConfig.AvgBlockTime = 1 * time.Second

	err = tmConfig.validate()
	require.Nil(t, err)

}

func TestValidateKafka(t *testing.T) {
	testKafkaConfig := kafkaConfig{}
	testKafkaConfig.ToTendermint.MinBatchSize = 1
	testKafkaConfig.ToEth.MinBatchSize = 1

	err := testKafkaConfig.validate()
	require.Equal(t, "replicationFactor has to be atleast 1", err.Error())
	testKafkaConfig.TopicDetail.ReplicationFactor = 1

	err = testKafkaConfig.validate()
	require.Equal(t, "num participants has to be atleast 1", err.Error())
	testKafkaConfig.TopicDetail.NumPartitions = 1

	err = testKafkaConfig.validate()
	require.Equal(t, "tendermint min batch size cannot be greater than max batch size", err.Error())
	testKafkaConfig.ToTendermint.MaxBatchSize = 2

	err = testKafkaConfig.validate()
	require.Equal(t, "ethereum min batch size cannot be greater than max batch size", err.Error())
	testKafkaConfig.ToEth.MaxBatchSize = 1

	err = testKafkaConfig.validate()
	require.Equal(t, "Kafka.MaxTendermintTxAttempts cannot be less than equal to 0", err.Error())
	testKafkaConfig.MaxTendermintTxAttempts = 5

	err = testKafkaConfig.validate()
	require.Equal(t, "kafka EthUnbondCycleTime time cannot be 0", err.Error())
	testKafkaConfig.EthUnbondCycleTime = 5 * time.Second

	err = testKafkaConfig.validate()
	require.Nil(t, err)
}

func TestValidateCASP(t *testing.T) {
	testConfig := caspConfig{}

	err := testConfig.validate()
	require.Equal(t, fmt.Errorf("casp vault id empty"), err)
	testConfig.VaultID = "vault id"

	err = testConfig.validate()
	require.Equal(t, fmt.Errorf("casp api token empty"), err)
	testConfig.ApiToken = "api token"

	err = testConfig.validate()
	require.Equal(t, fmt.Errorf("invalid casp url: parse \"\": empty url"), err)
	testConfig.URL = "https://127.0.0.1:443"

	err = testConfig.validate()
	require.Equal(t, fmt.Errorf("casp tendermint public empty"), err)
	testConfig.TendermintPublicKey = "TendermintPublicKey"

	err = testConfig.validate()
	require.Equal(t, fmt.Errorf("casp ethereum public empty"), err)
	testConfig.EthereumPublicKey = "EthereumPublicKey"

	err = testConfig.validate()
	require.Equal(t, fmt.Errorf("casp MaxAttempts cannot be equal to 0"), err)
	testConfig.MaxAttempts = 5

	err = testConfig.validate()
	require.Equal(t, fmt.Errorf("casp wait time cannot be 0"), err)
	testConfig.WaitTime = 5 * time.Second

	err = testConfig.validate()
	require.Nil(t, err)
}

func TestValidateTelegram(t *testing.T) {
	testTelegramConfig := newTelegramBot()
	err := testTelegramConfig.validate()
	require.Nil(t, err)

	testTelegramConfig.ChatID = 1
	err = testTelegramConfig.validate()
	require.Equal(t, "telegram bot configuration invalid", err.Error())
}

func TestConfig(t *testing.T) {
	teleConfig := telegramBot{
		Token:  "",
		ChatID: 1,
	}
	c := config{
		Kafka:       kafkaConfig{},
		Tendermint:  tendermintConfig{},
		Ethereum:    ethereumConfig{},
		CASP:        caspConfig{},
		TelegramBot: teleConfig,
		seal:        false,
		RPCEndpoint: "",
	}

	err := c.validate()
	require.Equal(t, fmt.Errorf("invalid EthereumEndPoint: parse \"\": empty url"), err)
	c.Ethereum.EthereumEndPoint = "ws://127.0.0.1:8546"
	c.Ethereum.GasLimit = 1
	c.Ethereum.GasFeeCap = 1
	c.Ethereum.LiquidStakingAddress = "0x0000000000000000000000000000000000000001"
	c.Ethereum.TokenWrapperAddress = "0x0000000000000000000000000000000000000001"
	c.Ethereum.bridgeAdminAddress = common.HexToAddress("0x0000000000000000000000000000000000000001")

	err = c.validate()
	require.Equal(t, fmt.Errorf("account prefix cannot be empty"), err)
	c.Tendermint.AccountPrefix = "cosmos"
	c.Tendermint.ChainID = "test"
	c.Tendermint.BroadcastMode = "sync"
	c.Tendermint.MinimumWrapAmount = 50
	c.Tendermint.GasAdjustment = 1.5
	c.Tendermint.GasPrice = "0.025"
	c.Tendermint.Node = "http://127.0.0.1:26657"
	c.Tendermint.Denom = "stake"
	c.Tendermint.CoinType = 118
	c.Tendermint.wrapAddress = "cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u"
	c.Tendermint.AvgBlockTime = 5 * time.Second

	err = c.validate()
	require.Equal(t, fmt.Errorf("replicationFactor has to be atleast 1"), err)
	c.Kafka.TopicDetail.ReplicationFactor = 1
	c.Kafka.TopicDetail.NumPartitions = 1
	c.Kafka.MaxTendermintTxAttempts = 5
	c.Kafka.EthUnbondCycleTime = 5 * time.Second

	err = c.validate()
	require.Equal(t, fmt.Errorf("casp vault id empty"), err)
	c.CASP.VaultID = "vault id"
	c.CASP.ApiToken = "api"
	c.CASP.URL = "https://127.0.0.1:443"
	c.CASP.MaxAttempts = 5
	c.CASP.EthereumPublicKey = "public key"
	c.CASP.TendermintPublicKey = "public key"
	c.CASP.WaitTime = 5 * time.Second

	err = c.validate()
	require.Equal(t, fmt.Errorf("telegram bot configuration invalid"), err)
	c.TelegramBot = telegramBot{}

	err = c.validate()
	require.Equal(t, fmt.Errorf("rpc endpoint empty"), err)
	c.RPCEndpoint = "http://localhost:4040"

	err = c.validate()
	require.Nil(t, err)
}
