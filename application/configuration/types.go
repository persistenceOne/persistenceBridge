package configuration

import (
	"crypto/ecdsa"
	"github.com/Shopify/sarama"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"time"
)

type Config struct {
	Kafka      KafkaConfig
	Tendermint TendermintConfig
	Ethereum   EthereumConfig
	PStakeHome string
}

func NewConfig() Config {
	return Config{
		Kafka:      NewKafkaConfig(),
		Tendermint: NewTendermintConfig(),
		Ethereum:   NewEthereumConfig(),
		PStakeHome: constants.DefaultPBridgeHome,
	}
}

type EthereumConfig struct {
	EthAccountPrivateKey *ecdsa.PrivateKey
	EthGasLimit          uint64
}

func NewEthereumConfig() EthereumConfig {
	return EthereumConfig{
		EthAccountPrivateKey: nil,
		EthGasLimit:          constants.DefaultEthGasLimit,
	}
}

type TendermintConfig struct {
	PStakeAddress sdkTypes.AccAddress
	PStakeDenom   string
}

func NewTendermintConfig() TendermintConfig {
	return TendermintConfig{
		PStakeAddress: nil,
		PStakeDenom:   constants.DefaultDenom,
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
