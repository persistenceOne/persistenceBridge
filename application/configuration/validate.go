package configuration

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"net/url"
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
	return nil
}

func (config ethereumConfig) validate() error {
	if config.GasLimit <= 0 {
		return fmt.Errorf("invalid eth gas limit")
	}
	if config.GasFeeCap <= 0 {
		return fmt.Errorf("invalid eth gas fee cap")
	}
	return nil
}

func (config tendermintConfig) validate() error {
	if config.pStakeAddress == "" {
		return fmt.Errorf("pStakeAddress empty")
	}
	_, err := sdk.AccAddressFromBech32(config.pStakeAddress)
	if err != nil {
		return err
	}
	if config.AccountPrefix == "" {
		return fmt.Errorf("account prefix cannot be empty")
	}
	if config.PStakeDenom == "" {
		return fmt.Errorf("denom cannot be empty")
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
	return nil
}

func (config caspConfig) validate() error {
	if config.VaultID == "" {
		return fmt.Errorf("casp vault id empty")
	}
	if config.apiToken == "" {
		return fmt.Errorf("casp api token empty")
	}
	if config.URL == "" {
		return fmt.Errorf("casp url empty")
	}
	if config.TendermintPublicKey == "" {
		return fmt.Errorf("casp tendermint public empty")
	}
	if config.EthereumPublicKey == "" {
		return fmt.Errorf("casp tendermint public empty")
	}
	if config.MaxAttempts <= 0 {
		return fmt.Errorf("casp MaxAttempts cannot be less than or equal to 0")
	}
	return nil
}

func (config telegramBot) validate() error {
	if (config.ChatID != 0 && config.Token == "") || (config.ChatID == 0 && config.Token != "") {
		return fmt.Errorf("telegram bot configuration invalid")
	}
	return nil
}
