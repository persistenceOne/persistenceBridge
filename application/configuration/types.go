package configuration

import (
	"crypto/ecdsa"
	"github.com/Shopify/sarama"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"time"
)

type Config struct {
	Kafka       KafkaConfig
	Tendermint  TendermintConfig
	Ethereum    EthereumConfig
	PStakeDenom string
	CoinType    uint32
	PStakeHome  string
}

func NewConfig() Config {
	return Config{
		Kafka:       NewKafkaConfig(),
		Tendermint:  NewTendermintConfig(),
		Ethereum:    NewEthereumConfig(),
		PStakeDenom: constants.DefaultDenom,
		CoinType:    constants.DefaultCoinType,
		PStakeHome:  constants.DefaultPBridgeHome,
	}
}

type EthereumConfig struct {
	EthAccountPrivateKey *ecdsa.PrivateKey
	EthGasLimit          uint64
	EthereumEndpoint     string
	EthereumSleepTime    int // seconds
	EthereumStartHeight  int64
}

func NewEthereumConfig() EthereumConfig {
	return EthereumConfig{
		EthAccountPrivateKey: nil,
		EthGasLimit:          constants.DefaultEthGasLimit,
		EthereumEndpoint:     constants.DefaultEthereumEndPoint,
		EthereumSleepTime:    constants.DefaultEthereumSleepTime,
		EthereumStartHeight:  constants.DefaultEthereumStartHeight,
	}
}

type TendermintConfig struct {
	PStakeAddress         sdkTypes.AccAddress
	Validators            []string
	RelayerTimeout        string
	TendermintSleepTime   int //seconds
	TendermintStartHeight int64
}

func NewTendermintConfig() TendermintConfig {
	return TendermintConfig{
		PStakeAddress:         nil,
		Validators:            []string{constants.Validator1.String()},
		RelayerTimeout:        constants.DefaultTimeout,
		TendermintSleepTime:   constants.DefaultTendermintSleepTime,
		TendermintStartHeight: constants.DefaultTendermintStartHeight,
	}
}

type KafkaConfig struct {
	// Brokers: List of brokers to run kafka cluster
	Brokers      []string
	TopicDetail  sarama.TopicDetail
	ToEth        TopicConsumer
	ToTendermint TopicConsumer
	// start the first unbond
	EthUnbondStartTime time.Duration
	// Time for each unbonding transactions 3 days => input nano-seconds 259200000000000
	EthUnbondCycleTime time.Duration
}

type TopicConsumer struct {
	MinBatchSize int
	MaxBatchSize int
	Ticker       time.Duration
}

func NewKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers:     constants.DefaultBrokers,
		TopicDetail: constants.TopicDetail,
		ToEth: TopicConsumer{
			MinBatchSize: constants.MinEthBatchSize,
			MaxBatchSize: constants.MaxEthBatchSize,
			Ticker:       constants.EthTicker,
		},
		ToTendermint: TopicConsumer{
			MinBatchSize: constants.MinTendermintBatchSize,
			MaxBatchSize: constants.MaxTendermintBatchSize,
			Ticker:       constants.TendermintTicker,
		},
		EthUnbondStartTime: time.Duration(time.Now().Unix()),
		EthUnbondCycleTime: constants.DefaultEthUnbondCycleTime,
	}
}
