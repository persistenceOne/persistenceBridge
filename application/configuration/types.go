package configuration

import (
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"time"
)

type Config struct {
	Kafka       KafkaConfig
	Tendermint  TendermintConfig
	Ethereum    EthereumConfig
	CASP        CASPConfig
	PBridgeHome string
	set         bool
	RPCEndpoint string
}

func NewConfig() Config {
	return Config{
		Kafka:       NewKafkaConfig(),
		Tendermint:  NewTendermintConfig(),
		Ethereum:    NewEthereumConfig(),
		CASP:        NewCASPConfig(),
		PBridgeHome: constants.DefaultPBridgeHome,
		set:         false,
		RPCEndpoint: constants.DefaultRPCEndpoint,
	}
}

type EthereumConfig struct {
	BridgeAdmin common.Address
	GasLimit    uint64
}

func NewEthereumConfig() EthereumConfig {
	return EthereumConfig{
		BridgeAdmin: constants.DefaultBridgeAdmin,
		GasLimit:    constants.DefaultEthGasLimit,
	}
}

type TendermintConfig struct {
	PStakeAddress string
	PStakeDenom   string
	BroadcastMode string
}

func NewTendermintConfig() TendermintConfig {
	return TendermintConfig{
		PStakeAddress: constants.DefaultPStakeAddress,
		PStakeDenom:   constants.DefaultDenom,
		BroadcastMode: constants.DefaultBroadcastMode,
	}
}

type CASPConfig struct {
	URL                     string
	VaultID                 string
	TendermintPublicKey     string
	EthereumPublicKey       string
	SignatureWaitTime       time.Duration
	APIToken                string
	AllowConcurrentKeyUsage bool
}

func NewCASPConfig() CASPConfig {
	return CASPConfig{
		URL:                     constants.DefaultCASPUrl,
		VaultID:                 constants.DefaultCASPVaultID,
		TendermintPublicKey:     constants.DefaultCASPTendermintPublicKey,
		EthereumPublicKey:       constants.DefaultCASPEthereumPublicKey,
		SignatureWaitTime:       constants.DefaultCASPSignatureWaitTime,
		APIToken:                constants.DefaultCASPAPI,
		AllowConcurrentKeyUsage: true,
	}
}

type KafkaConfig struct {
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
		EthUnbondCycleTime: constants.DefaultEthUnbondCycleTime,
	}
}
