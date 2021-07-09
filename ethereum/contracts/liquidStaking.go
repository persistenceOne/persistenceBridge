package contracts

import (
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
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
	address: constants2.LiquidStakingAddress,
	abi:     abi.ABI{},
	methods: map[string]func(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec, arguments []interface{}) error{
		constants2.LiquidStakingStake:   onStake,
		constants2.LiquidStakingUnStake: onUnStake,
	},
}

func onStake(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec, arguments []interface{}) error {
	amount := arguments[1].(*big.Int)
	stakeMsg := &stakingTypes.MsgDelegate{
		DelegatorAddress: configuration.GetAppConfig().Tendermint.PStakeAddress,
		ValidatorAddress: "",
		Amount:           sdkTypes.NewCoin(configuration.GetAppConfig().Tendermint.PStakeDenom, sdkTypes.NewInt(amount.Int64())),
	}

	msgBytes, err := protoCodec.MarshalInterface(sdkTypes.Msg(stakeMsg))
	if err != nil {
		log.Println("Failed to generate msgBytes: ", err)
		return err
	}
	log.Printf("Adding stake msg to kafka producer MsgDelegate: %s\n", stakeMsg.String())
	err = utils.ProducerDeliverMessage(msgBytes, utils.MsgDelegate, *kafkaProducer)
	if err != nil {
		log.Println("Failed to add msg to kafka queue: ", err)
		return err
	}
	return nil
}

func onUnStake(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec, arguments []interface{}) error {
	amount := arguments[1].(*big.Int)
	unStakeMsg := &stakingTypes.MsgUndelegate{
		DelegatorAddress: configuration.GetAppConfig().Tendermint.PStakeAddress,
		ValidatorAddress: "",
		Amount:           sdkTypes.NewCoin(configuration.GetAppConfig().Tendermint.PStakeDenom, sdkTypes.NewInt(amount.Int64())),
	}
	msgBytes, err := protoCodec.MarshalInterface(sdkTypes.Msg(unStakeMsg))
	if err != nil {
		log.Println("Failed to generate msgBytes: ", err)
		return err
	}
	log.Printf("Adding unStake msg to kafka producer EthUnbond: %s\n", unStakeMsg.String())
	err = utils.ProducerDeliverMessage(msgBytes, utils.EthUnbond, *kafkaProducer)
	if err != nil {
		log.Println("Failed to add msg to kafka queue: ", err)
		return err
	}
	return nil
}
