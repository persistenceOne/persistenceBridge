/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package configuration

import (
	"fmt"
	"net/url"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Validate: panics if config is not valid
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
func (cfg *ethereumConfig) validate() error {
	if cfg.GasLimit <= 0 {
		return ErrInvalidGasLimit
	}

	return nil
}

// Validate :panics if config is not valid
func (c *tendermintConfig) validate() error {
	if c.pStakeAddress == "" {
		return ErrPStakeAddressEmpty
	}

	_, err := sdk.AccAddressFromBech32(c.pStakeAddress)
	if err != nil {
		return err
	}

	if c.AccountPrefix == "" {
		return ErrEmptyAccountPrefix
	}

	if c.PStakeDenom == "" {
		return ErrEmptyDenom
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

	return nil
}

func (c *caspConfig) validate() error {
	if c.VaultID == "" {
		return ErrCaspVaultIDEmpty
	}

	if c.APIToken == "" {
		return ErrCaspAPITokenEmpty
	}

	if c.URL == "" {
		return ErrCaspURLEmpty
	}

	if c.TendermintPublicKey == "" {
		return fmt.Errorf("tendermint %w", ErrCaspPublicEmpty)
	}

	if c.EthereumPublicKey == "" {
		return fmt.Errorf("ethereum %w", ErrCaspPublicEmpty)
	}

	if c.MaxGetSignatureAttempts <= 0 {
		return ErrTooLowCaspMaxGetSignatureAttempts
	}

	return nil
}

func (c *telegramBot) validate() error {
	if (c.ChatID != 0 && c.Token == "") || (c.ChatID == 0 && c.Token != "") {
		return ErrTelegramBotInvalidConfig
	}

	return nil
}
