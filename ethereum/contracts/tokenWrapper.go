package contracts

import (
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"log"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
)

var TokenWrapper = Contract{
	name:    "TOKEN_WRAPPER",
	address: common.HexToAddress(constants2.TokenWrapperAddress),
	abi:     abi.ABI{},
	methods: map[string]func(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec, arguments []interface{}) error{
		constants2.TokenWrapperWithdrawUTokens: onWithdrawUTokens,
	},
}

func onWithdrawUTokens(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec, arguments []interface{}) error {
	ercAddress := arguments[0].(common.Address)
	amount := sdkTypes.NewIntFromBigInt(arguments[1].(*big.Int))
	atomAddress, err := sdkTypes.AccAddressFromBech32(arguments[2].(string))
	if err != nil {
		return err
	}
	sendCoinMsg := &bankTypes.MsgSend{
		FromAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
		ToAddress:   atomAddress.String(),
		Amount:      sdkTypes.NewCoins(sdkTypes.NewCoin(configuration.GetAppConfig().Tendermint.PStakeDenom, amount)),
	}
	msgBytes, err := protoCodec.MarshalInterface(sendCoinMsg)
	if err != nil {
		log.Println("Failed to generate msgBytes: ", err)
		return err
	}
	log.Printf("Adding sendCoin msg to kafka producer MsgSend, from: %s, to: %s, amount: %s\n", ercAddress.String(), atomAddress.String(), sendCoinMsg.Amount.String())
	err = utils.ProducerDeliverMessage(msgBytes, utils.MsgSend, *kafkaProducer)
	if err != nil {
		log.Printf("Failed to add msg to kafka queue MsgSend: %s [ETH Listener (onWithDrawUTokens)]\n", err.Error())
		return err
	}
	return nil
}
