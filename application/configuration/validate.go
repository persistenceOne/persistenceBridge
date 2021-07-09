package configuration

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Validate :panics if config is not valid
func (config Config) Validate() error {
	if err := config.Ethereum.Validate(); err != nil {
		panic(err)
	}
	if err := config.Tendermint.Validate(); err != nil {
		panic(err)
	}
	if err := config.Kafka.Validate(); err != nil {
		panic(err)
	}
	return nil
}

// Validate :panics if config is not valid
func (config EthereumConfig) Validate() error {

	return nil
}

// Validate :panics if config is not valid
func (config TendermintConfig) Validate() error {
	_, err := sdk.AccAddressFromBech32(config.PStakeAddress)
	if err != nil {
		return err
	}
	if !(config.BroadcastMode == flags.BroadcastAsync || config.BroadcastMode == flags.BroadcastSync || config.BroadcastMode == flags.BroadcastBlock) {
		return fmt.Errorf("invalid broadcast mode")
	}
	return nil
}

// Validate :panics if config is not valid
func (config KafkaConfig) Validate() error {
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
	return nil
}
