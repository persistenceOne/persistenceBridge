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
	GasFeeCap        int64
}

func newEthereumConfig() ethereumConfig {
	return ethereumConfig{
		EthereumEndPoint: constants.DefaultEthereumEndPoint,
		GasLimit:         constants.DefaultEthGasLimit,
		GasFeeCap:        constants.DefaultEthGasFeeCap,
	}
}

type tendermintConfig struct {
	pStakeAddress     string
	PStakeDenom       string
	BroadcastMode     string
	MinimumWrapAmount int64
	AccountPrefix     string
	Node              string
	ChainID           string
	CoinType          uint32
}

func newTendermintConfig() tendermintConfig {
	return tendermintConfig{
		PStakeDenom:       constants.DefaultDenom,
		BroadcastMode:     constants.DefaultBroadcastMode,
		MinimumWrapAmount: constants.DefaultMinimumWrapAmount,
		AccountPrefix:     constants.DefaultAccountPrefix,
		Node:              constants.DefaultTendermintNode,
		ChainID:           constants.DefaultTendermintChainId,
		CoinType:          constants.DefaultTendermintCoinType,
	}
}

type caspConfig struct {
	URL                     string
	VaultID                 string
	TendermintPublicKey     string
	EthereumPublicKey       string
	WaitTime                time.Duration
	APIToken                string
	AllowConcurrentKeyUsage bool
	MaxAttempts             int
}

func newCASPConfig() caspConfig {
	return caspConfig{
		URL:                     "",
		VaultID:                 "",
		TendermintPublicKey:     "",
		EthereumPublicKey:       "",
		WaitTime:                constants.DefaultCASPWaitTime,
		APIToken:                "",
		AllowConcurrentKeyUsage: true,
		MaxAttempts:             constants.DefaultCASPMaxAttempts,
	}
}

type kafkaConfig struct {
	Brokers                 []string // Brokers: List of brokers to run kafka cluster
	TopicDetail             sarama.TopicDetail
	ToEth                   TopicConsumer
	ToTendermint            TopicConsumer
	EthUnbondCycleTime      time.Duration // Time for each unbonding transactions 3 days => input nano-seconds 259200000000000
	MaxTendermintTxAttempts int           // Max attempts in kafka consumer toTendermint to do tx if there is already a tx
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
		EthUnbondCycleTime:      constants.DefaultEthUnbondCycleTime,
		MaxTendermintTxAttempts: constants.DefaultTendermintMaxTxAttempts,
	}
}

func (config tendermintConfig) GetPStakeAddress() string {
	if config.pStakeAddress == "" {
		log.Fatalln("pStakeAddress not set")
	}
	return config.pStakeAddress
}
