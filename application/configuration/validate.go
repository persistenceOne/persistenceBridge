/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"net/url"
	"strconv"
)

func (config config) validate() error {
	if err := config.Ethereum.validate(); err != nil {
		return err
	}
	if err := config.Tendermint.validate(); err != nil {
		return err
	}
	if err := config.Kafka.validate(); err != nil {
		return err
	}
	if err := config.CASP.validate(); err != nil {
		return err
	}
	if err := config.TelegramBot.validate(); err != nil {
		return err
	}
	if config.RPCEndpoint == "" {
		return fmt.Errorf("rpc endpoint empty")
	}
	if config.InitSlackBot == false {
		return fmt.Errorf("slack alerts not set")
	}
	if config.SlackBotToken == "" {
		return fmt.Errorf("slack bot configuration invalid")
	}
	return nil
}

func (config ethereumConfig) validate() error {
	if _, err := url.ParseRequestURI(config.EthereumEndPoint); err != nil {
		return fmt.Errorf("invalid EthereumEndPoint: %v", err)
	}
	if config.GasLimit == 0 {
		return fmt.Errorf("invalid eth gas limit")
	}
	if config.GasFeeCap <= 0 {
		return fmt.Errorf("invalid eth gas fee cap")
	}
	if config.LiquidStakingAddress == constants.EthereumZeroAddress || !common.IsHexAddress(config.LiquidStakingAddress) {
		return fmt.Errorf("empty liquid staking contract address")
	}
	if config.TokenWrapperAddress == constants.EthereumZeroAddress || !common.IsHexAddress(config.TokenWrapperAddress) {
		return fmt.Errorf("empty token wrapper contract address")
	}
	if config.bridgeAdminAddress.String() == constants.EthereumZeroAddress {
		return fmt.Errorf("bridgeAdminAddress is empty")
	}
	if config.BalanceCheckPeriod <= 0 || config.AlertAmount <= 0 {
		return fmt.Errorf("invalid ethereum balance alert configuration")
	}
	return nil
}

func (config tendermintConfig) validate() error {
	if config.AccountPrefix == "" {
		return fmt.Errorf("account prefix cannot be empty")
	}
	if config.Denom == "" {
		return fmt.Errorf("denom cannot be empty")
	}
	if _, err := strconv.ParseFloat(config.GasPrice, 64); err != nil {
		return fmt.Errorf("invalied tendermint gas price %v", err)
	}
	if config.GasAdjustment <= 1.0 {
		return fmt.Errorf("tendermint gas adjustment should be greater than 1 (recommended 1.5, current: %v)", config.GasAdjustment)
	}
	if config.MinimumWrapAmount < 0 {
		return fmt.Errorf("minimum wrap amount cannot be less than 0")
	}
	if config.ChainID == "" {
		return fmt.Errorf("chain id cannot be empty")
	}
	if _, err := url.ParseRequestURI(config.Node); err != nil {
		return fmt.Errorf("invalid tendermint node: %v", err)
	}
	if !(config.BroadcastMode == flags.BroadcastAsync || config.BroadcastMode == flags.BroadcastSync || config.BroadcastMode == flags.BroadcastBlock) {
		return fmt.Errorf("invalid broadcast mode")
	}
	if config.wrapAddress == "" {
		return fmt.Errorf("wrapAddress empty")
	}
	_, err := sdk.AccAddressFromBech32(config.wrapAddress)
	if err != nil {
		return err
	}
	if config.AvgBlockTime.Nanoseconds() == 0 {
		return fmt.Errorf("tendermint chain avg block time cannot be 0")
	}
	return nil
}

// Validate :panics if config is not valid
func (config kafkaConfig) validate() error {
	if config.TopicDetail.ReplicationFactor < 1 {
		return errors.New("replicationFactor has to be atleast 1")
	}
	if config.TopicDetail.NumPartitions < 1 {
		return errors.New("num participants has to be atleast 1")
	}
	if config.ToTendermint.MinBatchSize > config.ToTendermint.MaxBatchSize {
		return errors.New("tendermint min batch size cannot be greater than max batch size")
	}
	if config.ToEth.MinBatchSize > config.ToEth.MaxBatchSize {
		return errors.New("ethereum min batch size cannot be greater than max batch size")
	}
	if config.MaxTendermintTxAttempts <= 0 {
		return errors.New("Kafka.MaxTendermintTxAttempts cannot be less than equal to 0")
	}
	if config.EthUnbondCycleTime.Nanoseconds() == 0 {
		return fmt.Errorf("kafka EthUnbondCycleTime time cannot be 0")
	}
	return nil
}

func (config caspConfig) validate() error {
	if config.VaultID == "" {
		return fmt.Errorf("casp vault id empty")
	}
	if config.ApiToken == "" {
		return fmt.Errorf("casp api token empty")
	}
	if _, err := url.ParseRequestURI(config.URL); err != nil {
		return fmt.Errorf("invalid casp url: %v", err)
	}
	if config.TendermintPublicKey == "" {
		return fmt.Errorf("casp tendermint public empty")
	}
	if config.EthereumPublicKey == "" {
		return fmt.Errorf("casp ethereum public empty")
	}
	if config.MaxAttempts == 0 {
		return fmt.Errorf("casp MaxAttempts cannot be equal to 0")
	}
	if config.WaitTime.Nanoseconds() == 0 {
		return fmt.Errorf("casp wait time cannot be 0")
	}
	return nil
}

func (config telegramBot) validate() error {
	if (config.ChatID != 0 && config.Token == "") || (config.ChatID == 0 && config.Token != "") {
		return fmt.Errorf("telegram bot configuration invalid")
	}
	return nil
}


