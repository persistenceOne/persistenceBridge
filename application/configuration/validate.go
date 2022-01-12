/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/constants"
)

func (c config) validate() error {
	if err := c.Ethereum.validate(); err != nil {
		return err
	}

	if err := c.Tendermint.validate(); err != nil {
		return err
	}

	if err := c.Kafka.validate(); err != nil {
		return err
	}

	if err := c.CASP.validate(); err != nil {
		return err
	}

	if err := c.TelegramBot.validate(); err != nil {
		return err
	}

	if c.RPCEndpoint == "" {
		return ErrRPCEndpointEmpty
	}

	return nil
}

// Validate :panics if config is not valid
func (c ethereumConfig) validate() error {
	if _, err := url.ParseRequestURI(c.EthereumEndPoint); err != nil {
		return fmt.Errorf("invalid EthereumEndPoint: %v", err)
	}

	if c.GasLimit == 0 {
		return ErrInvalidGasLimit
	}

	if c.GasFeeCap <= 0 {
		return fmt.Errorf("invalid eth gas fee cap")
	}

	if c.LiquidStakingAddress == constants.EthereumZeroAddress() {
		return fmt.Errorf("empty liquid staking contract address")
	}

	if c.TokenWrapperAddress == constants.EthereumZeroAddress() {
		return fmt.Errorf("empty token wrapper contract address")
	}

	if c.bridgeAdminAddress == constants.EthereumZeroAddress() {
		return fmt.Errorf("bridgeAdminAddress is empty")
	}

	if c.BalanceCheckPeriod <= 0 || c.AlertAmount <= 0 {
		return fmt.Errorf("invalid ethereum balance alert configuration")
	}

	return nil
}

func (c tendermintConfig) validate() error {
	if c.AccountPrefix == "" {
		return fmt.Errorf("account prefix cannot be empty")
	}

	if c.Denom == "" {
		return fmt.Errorf("denom cannot be empty")
	}

	if _, err := strconv.ParseFloat(c.GasPrice, 64); err != nil {
		return fmt.Errorf("invalied tendermint gas price %v", err)
	}

	// fixme: float-point compare
	if c.GasAdjustment <= 1.0 {
		return fmt.Errorf("tendermint gas adjustment should be greater than 1 (recommended 1.5, current: %v)", c.GasAdjustment)
	}

	if c.MinimumWrapAmount < 0 {
		return ErrNegativeWrapAmount
	}

	if c.ChainID == "" {
		return ErrEmptyChainID
	}

	if _, err := url.ParseRequestURI(c.Node); err != nil {
		// nolint it already has %w
		// nolint: errorlint
		return fmt.Errorf("%w: %s", ErrInvalidTendermintNode, err.Error())
	}

	if !(c.BroadcastMode == flags.BroadcastAsync || c.BroadcastMode == flags.BroadcastSync || c.BroadcastMode == flags.BroadcastBlock) {
		return ErrInvalidBroadcastMode
	}

	if c.wrapAddress == "" {
		return fmt.Errorf("wrapAddress empty")
	}

	_, err := sdk.AccAddressFromBech32(c.wrapAddress)
	if err != nil {
		return err
	}

	if c.AvgBlockTime.Nanoseconds() == 0 {
		return fmt.Errorf("tendermint chain avg block time cannot be 0")
	}

	return nil
}

// Validate :panics if config is not valid
func (c *kafkaConfig) validate() error {
	if c.TopicDetails.ReplicationFactor < 1 {
		return ErrTooLowReplicationFactor
	}

	if c.TopicDetails.NumPartitions < 1 {
		return ErrTooFewParticipants
	}

	if c.ToTendermint.MinBatchSize > c.ToTendermint.MaxBatchSize {
		return fmt.Errorf("tendermint %w", ErrTooBigMinBatchSize)
	}

	if c.ToEth.MinBatchSize > c.ToEth.MaxBatchSize {
		return fmt.Errorf("ethereum %w", ErrTooBigMinBatchSize)
	}

	if c.MaxTendermintTxAttempts <= 0 {
		return errors.New("Kafka.MaxTendermintTxAttempts cannot be less than equal to 0")
	}

	if c.EthUnbondCycleTime.Nanoseconds() == 0 {
		return fmt.Errorf("kafka EthUnbondCycleTime time cannot be 0")
	}

	return nil
}

func (c *caspConfig) validate() error {
	if c.VaultID == "" {
		return ErrCaspVaultIDEmpty
	}

	if c.APIToken == "" {
		return ErrCaspAPITokenEmpty
	}

	if _, err := url.ParseRequestURI(c.URL); err != nil {
		return fmt.Errorf("invalid casp url: %v", err)
	}

	if c.TendermintPublicKey == "" {
		return fmt.Errorf("tendermint %w", ErrCaspPublicEmpty)
	}

	if c.EthereumPublicKey == "" {
		return fmt.Errorf("ethereum %w", ErrCaspPublicEmpty)
	}

	if c.MaxAttempts == 0 {
		return fmt.Errorf("casp MaxAttempts cannot be equal to 0")
	}

	if c.WaitTime.Nanoseconds() == 0 {
		return fmt.Errorf("casp wait time cannot be 0")
	}

	return nil
}

func (c *telegramBot) validate() error {
	if (c.ChatID != 0 && c.Token == "") || (c.ChatID == 0 && c.Token != "") {
		return ErrTelegramBotInvalidConfig
	}

	return nil
}
