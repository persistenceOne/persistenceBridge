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
func (config *config) validate() error {
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
		return ErrRPCEndpointEmpty
	}

	return nil
}

// Validate :panics if config is not valid
func (config *ethereumConfig) validate() error {
	if config.GasLimit <= 0 {
		return ErrInvalidGasLimit
	}

	return nil
}

// Validate :panics if config is not valid
func (config *tendermintConfig) validate() error {
	if config.pStakeAddress == "" {
		return ErrPStakeAddressEmpty
	}

	_, err := sdk.AccAddressFromBech32(config.pStakeAddress)
	if err != nil {
		return err
	}

	if config.AccountPrefix == "" {
		return ErrEmptyAccountPrefix
	}

	if config.PStakeDenom == "" {
		return ErrEmptyDenom
	}

	if config.MinimumWrapAmount < 0 {
		return ErrNegativeWrapAmount
	}

	if config.ChainID == "" {
		return ErrEmptyChainID
	}

	if _, err := url.ParseRequestURI(config.Node); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidTendermintNode, err)
	}

	if !(config.BroadcastMode == flags.BroadcastAsync || config.BroadcastMode == flags.BroadcastSync || config.BroadcastMode == flags.BroadcastBlock) {
		return ErrInvalidBroadcastMode
	}

	return nil
}

// Validate :panics if config is not valid
func (config *kafkaConfig) validate() error {
	if config.TopicDetail.ReplicationFactor < 1 {
		return ErrTooLowReplicationFactor
	}

	if config.TopicDetail.NumPartitions < 1 {
		return ErrTooFewParticipants
	}

	if config.ToTendermint.MinBatchSize > config.ToTendermint.MaxBatchSize {
		return fmt.Errorf("tendermint %w", ErrTooBigMinBatchSize)
	}

	if config.ToEth.MinBatchSize > config.ToEth.MaxBatchSize {
		return fmt.Errorf("ethereum %w", ErrTooBigMinBatchSize)
	}

	return nil
}

func (config *caspConfig) validate() error {
	if config.VaultID == "" {
		return ErrCaspVaultIDEmpty
	}

	if config.APIToken == "" {
		return ErrCaspAPITokenEmpty
	}

	if config.URL == "" {
		return ErrCaspURLEmpty
	}

	if config.TendermintPublicKey == "" {
		return fmt.Errorf("tendermint %w", ErrCaspPublicEmpty)
	}

	if config.EthereumPublicKey == "" {
		return fmt.Errorf("ethereum %w", ErrCaspPublicEmpty)
	}

	if config.MaxGetSignatureAttempts <= 0 {
		return ErrTooLowCaspMaxGetSignatureAttempts
	}

	return nil
}

func (config *telegramBot) validate() error {
	if (config.ChatID != 0 && config.Token == "") || (config.ChatID == 0 && config.Token != "") {
		return ErrTelegramBotInvalidConfig
	}

	return nil
}
