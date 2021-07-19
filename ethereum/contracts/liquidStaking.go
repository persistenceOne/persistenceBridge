package contracts

import (
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	"math/big"

	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var LiquidStaking = Contract{
	name:    "LIQUID_STAKING",
	address: common.HexToAddress(constants2.LiquidStakingAddress),
	abi:     abi.ABI{},
	methods: map[string]func(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec, arguments []interface{}) error{
		constants2.LiquidStakingStake:   onStake,
		constants2.LiquidStakingUnStake: onUnStake,
	},
}

func onStake(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec, arguments []interface{}) error {
	ercAddress := arguments[0].(common.Address)
	amount := sdkTypes.NewIntFromBigInt(arguments[1].(*big.Int))
	stakeMsg := &stakingTypes.MsgDelegate{
		DelegatorAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
		ValidatorAddress: "",
		Amount:           sdkTypes.NewCoin(configuration.GetAppConfig().Tendermint.PStakeDenom, amount),
	}

	msgBytes, err := protoCodec.MarshalInterface(stakeMsg)
	if err != nil {
		return err
	}
	logging.Info("Adding stake msg to kafka producer MsgDelegate, from:", ercAddress.String(), "amount:", amount.String())
	err = utils.ProducerDeliverMessage(msgBytes, utils.MsgDelegate, *kafkaProducer)
	if err != nil {
		logging.Error("Failed to add msg to kafka queue [ETH Listener (onStake)] MsgDelegate:", err)
		return err
	}
	return nil
}

func onUnStake(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec, arguments []interface{}) error {
	ercAddress := arguments[0].(common.Address)
	amount := sdkTypes.NewIntFromBigInt(arguments[1].(*big.Int))
	unStakeMsg := &stakingTypes.MsgUndelegate{
		DelegatorAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
		ValidatorAddress: "",
		Amount:           sdkTypes.NewCoin(configuration.GetAppConfig().Tendermint.PStakeDenom, amount),
	}
	msgBytes, err := protoCodec.MarshalInterface(unStakeMsg)
	if err != nil {
		return err
	}
	logging.Info("Adding unStake msg to kafka producer EthUnbond, from: ", ercAddress.String(), "amount:", amount.String())
	err = utils.ProducerDeliverMessage(msgBytes, utils.EthUnbond, *kafkaProducer)
	if err != nil {
		logging.Error("Failed to add msg to kafka queue  [ETH Listener (onUnStake)] EthUnbond:", err)
		return err
	}
	return nil
}
