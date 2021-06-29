package contracts

import (
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
	methods: map[string]func(kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec, arguments []interface{}) error{
		constants2.LiquidStakingStake:   onStake,
		constants2.LiquidStakingUnStake: onUnStake,
	},
}

func onStake(kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec, arguments []interface{}) error {
	amount := arguments[1].(*big.Int)
	stakeMsg := stakingTypes.NewMsgDelegate(configuration.GetAppConfiguration().PStakeAddress, constants2.Validator1, sdkTypes.NewCoin(configuration.GetAppConfiguration().PStakeDenom, sdkTypes.NewInt(amount.Int64())))
	msgBytes, err := protoCodec.MarshalInterface(sdkTypes.Msg(stakeMsg))
	if err != nil {
		log.Print("Failed to generate msgBytes: ", err)
		return err
	}
	log.Printf("Adding stake msg to kafka producer MsgDelegate: %s\n", stakeMsg.String())
	err = utils.ProducerDeliverMessage(msgBytes, utils.MsgDelegate, kafkaState.Producer)
	if err != nil {
		log.Print("Failed to add msg to kafka queue: ", err)
		return err
	}
	return nil
}

func onUnStake(kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec, arguments []interface{}) error {
	amount := arguments[1].(*big.Int)
	unStakeMsg := stakingTypes.NewMsgUndelegate(configuration.GetAppConfiguration().PStakeAddress, constants2.Validator1, sdkTypes.NewCoin(configuration.GetAppConfiguration().PStakeDenom, sdkTypes.NewInt(amount.Int64())))
	msgBytes, err := protoCodec.MarshalInterface(sdkTypes.Msg(unStakeMsg))
	if err != nil {
		log.Print("Failed to generate msgBytes: ", err)
		return err
	}
	log.Printf("Adding unStake msg to kafka producer EthUnbond: %s\n", unStakeMsg.String())
	err = utils.ProducerDeliverMessage(msgBytes, utils.EthUnbond, kafkaState.Producer)
	if err != nil {
		log.Print("Failed to add msg to kafka queue: ", err)
		return err
	}
	return nil
}
