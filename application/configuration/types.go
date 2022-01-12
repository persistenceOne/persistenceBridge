/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"fmt"
	"log"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/persistenceOne/persistenceBridge/application/constants"
)

//go:generate deep-copy --type config -o ./types_copy_gen.go ./types.go

type config struct {
	Kafka       kafkaConfig
	Tendermint  tendermintConfig
	Ethereum    ethereumConfig
	CASP        caspConfig
	TelegramBot telegramBot
	seal        bool
	RPCEndpoint string
	isFilled    bool
}

func (c config) IsSealed() bool {
	return c.seal
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
	EthereumEndPoint     string
	GasLimit             uint64
	GasFeeCap            int64
	bridgeAdminAddress   common.Address
	TokenWrapperAddress  common.Address
	LiquidStakingAddress common.Address
	BalanceCheckPeriod   uint64
	AlertAmount          int64
}

func newEthereumConfig() ethereumConfig {
	return ethereumConfig{
		EthereumEndPoint: constants.DefaultEthereumEndPoint,
		GasLimit:         constants.DefaultEthGasLimit,
		GasFeeCap:        constants.DefaultEthGasFeeCap,
	}
}

type tendermintConfig struct {
	wrapAddress       string
	Denom             string
	BroadcastMode     string
	GasPrice          string
	GasAdjustment     float64
	MinimumWrapAmount int64
	AccountPrefix     string
	Node              string
	ChainID           string
	CoinType          uint32
	AvgBlockTime      time.Duration
}

func newTendermintConfig() tendermintConfig {
	return tendermintConfig{
		Denom:             constants.DefaultDenom,
		BroadcastMode:     constants.DefaultBroadcastMode,
		GasPrice:          constants.DefaultTendermintGasPrice,
		GasAdjustment:     constants.DefaultTendermintGasAdjustment,
		MinimumWrapAmount: constants.DefaultMinimumWrapAmount,
		AccountPrefix:     constants.DefaultAccountPrefix,
		Node:              constants.DefaultTendermintNode,
		ChainID:           constants.DefaultTendermintChainID,
		CoinType:          constants.DefaultTendermintCoinType,
		AvgBlockTime:      constants.DefaultTendermintAvgBlockTime,
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
	MaxAttempts             uint
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
	// Brokers: List of brokers to run kafka cluster
	Brokers []string
	constants.TopicDetails
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
		Brokers:      []string{constants.DefaultBroker},
		TopicDetails: constants.TopicDetail,
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

func (config tendermintConfig) GetWrapAddress() string {
	if config.wrapAddress == "" {
		log.Fatalln("wrapAddress not set")
	}
	return config.wrapAddress
}

func setWrapAddress(tmAddress sdk.AccAddress) {
	if !appConfig.seal {
		if strings.Contains(tmAddress.String(), GetAppConfig().Tendermint.AccountPrefix) {
			appConfig.Tendermint.wrapAddress = tmAddress.String()
		} else {
			panic(fmt.Errorf("pStake wrap address prefix (%s) and Config account prefix (%s) does not match", sdk.GetConfig().GetBech32AccountAddrPrefix(), GetAppConfig().Tendermint.AccountPrefix))
		}
	}
}

func (c tendermintConfig) GetPStakeAddress() string {
	if c.wrapAddress == "" {
		log.Fatalln("pStakeAddress not set")
	}

	return c.wrapAddress
}

func setBridgeAdminAddress(address common.Address) {
	if appConfig.seal {
		return
	}

	if address == constants.EthereumZeroAddress() {
		panic(fmt.Errorf("invalid eth address"))
	}

	appConfig.Ethereum.bridgeAdminAddress = address
}

func (c ethereumConfig) GetBridgeAdminAddress() common.Address {
	return c.bridgeAdminAddress
}
