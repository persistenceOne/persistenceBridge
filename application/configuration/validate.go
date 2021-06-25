package configuration

import (
	"errors"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
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
	if config.EthereumStartHeight < 0 {
		return errors.New("EthereumStartHeight cannot be less than 0")
	}
	if config.EthereumSleepTime < 0 {
		return errors.New("EthereumSleepTime cannot be less than 0")
	}
	return nil
}

// Validate :panics if config is not valid
func (config TendermintConfig) Validate() error {
	for _, validator := range config.Validators {
		_, err := sdkTypes.ValAddressFromBech32(validator)
		if err != nil {
			return err
		}
	}
	if config.TendermintStartHeight < 0 {
		return errors.New("TendermintStartHeight cannot be less than 0")
	}
	if config.TendermintSleepTime < 0 {
		return errors.New("TendermintSleepTime cannot be less than 0")
	}
	return nil
}

// Validate :panics if config is not valid
func (config KafkaConfig) Validate() error {
	if config.TopicDetail.ReplicationFactor < 1 {
		errors.New("replicationFactor has to be atleast 1")
	}
	if config.TopicDetail.NumPartitions < 1 {
		errors.New("num participants has to be atleast 1")
	}
	if config.ToTendermint.MinBatchSize > config.ToTendermint.MaxBatchSize {
		errors.New("tendermint min batch size cannot be greater than max batch size")
	}
	if config.ToEth.MinBatchSize > config.ToEth.MaxBatchSize {
		errors.New("ethereum min batch size cannot be greater than max batch size")
	}
	return nil
}
