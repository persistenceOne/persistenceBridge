package configuration

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/persistenceOne/persistenceBridge/application/constants"
)

type config struct {
	Kafka       kafkaConfig
	Tendermint  tendermintConfig
	Ethereum    ethereumConfig
	CASP        caspConfig
	TelegramBot telegramBot
	seal        bool
	RPCEndpoint string
}

func newConfig() config {
	return config{
		Kafka:       newKafkaConfig(),
		Tendermint:  newTendermintConfig(),
		Ethereum:    newEthereumConfig(),
		CASP:        newCASPConfig(),
		TelegramBot: newTelegramBot(),
		seal:        false,
		RPCEndpoint: constants.DefaultRPCEndpoint,
	}
}

type ethereumConfig struct {
	EthereumEndPoint string
	GasLimit         uint64
}

func newEthereumConfig() ethereumConfig {
	return ethereumConfig{
		EthereumEndPoint: constants.DefaultEthereumEndPoint,
		GasLimit:         constants.DefaultEthGasLimit,
	}
}

type tendermintConfig struct {
	pStakeAddress     string
	PStakeDenom       string
	BroadcastMode     string
	MinimumWrapAmount int64
}

func newTendermintConfig() tendermintConfig {
	return tendermintConfig{
		PStakeDenom:       constants.DefaultDenom,
		BroadcastMode:     constants.DefaultBroadcastMode,
		MinimumWrapAmount: constants.DefaultMinimumWrapAmount,
	}
}

type caspConfig struct {
	URL                     string
	VaultID                 string
	TendermintPublicKey     string
	EthereumPublicKey       string
	SignatureWaitTime       time.Duration
	APIToken                string
	AllowConcurrentKeyUsage bool
	MaxGetSignatureAttempts int
}

func newCASPConfig() caspConfig {
	return caspConfig{
		URL:                     "",
		VaultID:                 "",
		TendermintPublicKey:     "",
		EthereumPublicKey:       "",
		SignatureWaitTime:       constants.DefaultCASPSignatureWaitTime,
		APIToken:                "",
		AllowConcurrentKeyUsage: true,
		MaxGetSignatureAttempts: constants.DefaultCASPMaxGetSignatureAttempt,
	}
}

type kafkaConfig struct {
	// Brokers: List of brokers to run kafka cluster
	Brokers      []string
	TopicDetail  sarama.TopicDetail
	ToEth        TopicConsumer
	ToTendermint TopicConsumer
	// Time for each unbonding transactions 3 days => input nano-seconds 259200000000000
	EthUnbondCycleTime time.Duration
}

type TopicConsumer struct {
	MinBatchSize int
	MaxBatchSize int
	Ticker       time.Duration
}

type telegramBot struct {
	Token  string
	ChatID int64
}

func newTelegramBot() telegramBot {
	return telegramBot{
		Token:  "",
		ChatID: 0,
	}
}

func newKafkaConfig() kafkaConfig {
	return kafkaConfig{
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
		EthUnbondCycleTime: constants.DefaultEthUnbondCycleTime,
	}
}

func (config tendermintConfig) GetPStakeAddress() string {
	if config.pStakeAddress == "" {
		log.Fatalln("pStakeAddress not set")
	}
	return config.pStakeAddress
}
