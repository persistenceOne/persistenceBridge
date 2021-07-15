package contracts

import (
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
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
	amount := sdkTypes.NewIntFromBigInt(arguments[1].(*big.Int))
	stakeMsg := &stakingTypes.MsgDelegate{
		DelegatorAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
		ValidatorAddress: "",
		Amount:           sdkTypes.NewCoin(configuration.GetAppConfig().Tendermint.PStakeDenom, amount),
	}

	msgBytes, err := protoCodec.MarshalInterface(stakeMsg)
	if err != nil {
		log.Println("Failed to generate msgBytes: ", err)
		return err
	}
	log.Printf("Adding stake msg to kafka producer MsgDelegate, amount: %s\n", amount.String())
	err = utils.ProducerDeliverMessage(msgBytes, utils.MsgDelegate, *kafkaProducer)
	if err != nil {
		log.Printf("Failed to add msg to kafka queue MsgDelegate: %s [ETH Listener (onStake)]\n", err.Error())
		return err
	}
	return nil
}

func onUnStake(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec, arguments []interface{}) error {
	amount := sdkTypes.NewIntFromBigInt(arguments[1].(*big.Int))
	unStakeMsg := &stakingTypes.MsgUndelegate{
		DelegatorAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
		ValidatorAddress: "",
		Amount:           sdkTypes.NewCoin(configuration.GetAppConfig().Tendermint.PStakeDenom, amount),
	}
	msgBytes, err := protoCodec.MarshalInterface(unStakeMsg)
	if err != nil {
		log.Println("Failed to generate msgBytes: ", err)
		return err
	}
	log.Printf("Adding unstake msg to kafka producer EthUnbond, amount: %s\n", amount.String())
	err = utils.ProducerDeliverMessage(msgBytes, utils.EthUnbond, *kafkaProducer)
	if err != nil {
		log.Printf("Failed to add msg to kafka queue EthUnbond: %s [ETH Listener (onUnStake)]\n", err.Error())
		return err
	}
	return nil
}
