package configuration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateEthereum(t *testing.T) {
	testEthConfig := newEthereumConfig()
	err := testEthConfig.validate()
	require.Nil(t, err)
	testEthConfig.GasLimit = 0
	err = testEthConfig.validate()
	require.Equal(t, "invalid eth gas limit", err.Error())
}
func TestValidateTendermint(t *testing.T) {
	testConfig := InitConfig()
	pstakeAddress, err := sdk.AccAddressFromBech32("cosmos1lfeqaqld74e2mmatx8luut0r4fajfu7kh3580u")
	require.Nil(t, err)
	testConfig.seal = false
	SetPStakeAddress(pstakeAddress)
	err = testConfig.Tendermint.validate()
	require.Nil(t, err)
	testConfig.Tendermint.pStakeAddress = ""
	err = testConfig.Tendermint.validate()
	require.Equal(t, "pStakeAddress empty", err.Error())

	testConfig.Tendermint = newTendermintConfig()
	SetPStakeAddress(pstakeAddress)
	testConfig.Tendermint.AccountPrefix = ""
	err = testConfig.Tendermint.validate()
	require.Equal(t, "account prefix cannot be empty", err.Error())

	testConfig.Tendermint = newTendermintConfig()
	SetPStakeAddress(pstakeAddress)
	testConfig.Tendermint.PStakeDenom = ""
	err = testConfig.Tendermint.validate()
	require.Equal(t, "denom cannot be empty", err.Error())

	testConfig.Tendermint = newTendermintConfig()
	SetPStakeAddress(pstakeAddress)
	testConfig.Tendermint.MinimumWrapAmount = -1
	err = testConfig.Tendermint.validate()
	require.Equal(t, "minimum wrap amount cannot be less than 0", err.Error())

	testConfig.Tendermint = newTendermintConfig()
	SetPStakeAddress(pstakeAddress)
	testConfig.Tendermint.ChainID = ""
	err = testConfig.Tendermint.validate()
	require.Equal(t, "chain id cannot be empty", err.Error())

	testConfig.Tendermint = newTendermintConfig()
	SetPStakeAddress(pstakeAddress)
	testConfig.Tendermint.Node = ""
	err = testConfig.Tendermint.validate()
	require.NotNil(t, err)

	testConfig.Tendermint = newTendermintConfig()
	SetPStakeAddress(pstakeAddress)
	testConfig.Tendermint.BroadcastMode = ""
	err = testConfig.Tendermint.validate()
	require.Equal(t, "invalid broadcast mode", err.Error())
}

func TestValidateKafka(t *testing.T) {
	testKafkaConfig := newKafkaConfig()
	err := testKafkaConfig.validate()
	require.Nil(t, err)

	testKafkaConfig.TopicDetail.ReplicationFactor = 0
	err = testKafkaConfig.validate()
	require.Equal(t, "replicationFactor has to be atleast 1", err.Error())

	testKafkaConfig = newKafkaConfig()
	testKafkaConfig.TopicDetail.NumPartitions = 0
	err = testKafkaConfig.validate()
	require.Equal(t, "num participants has to be atleast 1", err.Error())

	testKafkaConfig = newKafkaConfig()
	testKafkaConfig.ToTendermint.MinBatchSize = testKafkaConfig.ToTendermint.MaxBatchSize + 1
	err = testKafkaConfig.validate()
	require.Equal(t, "tendermint min batch size cannot be greater than max batch size", err.Error())

	testKafkaConfig = newKafkaConfig()
	testKafkaConfig.ToEth.MinBatchSize = testKafkaConfig.ToEth.MaxBatchSize + 1
	err = testKafkaConfig.validate()
	require.Equal(t, "ethereum min batch size cannot be greater than max batch size", err.Error())
}

func TestValidateCASP(t *testing.T) {
	testConfig := InitConfig()
	err := testConfig.validate()
	require.NotNil(t, err)
}

func TestValidateTelegram(t *testing.T) {
	testTelegramConfig := newTelegramBot()
	err := testTelegramConfig.validate()
	require.Nil(t, err)

	testTelegramConfig.ChatID = 1
	err = testTelegramConfig.validate()
	require.Equal(t, "telegram bot configuration invalid", err.Error())
}

