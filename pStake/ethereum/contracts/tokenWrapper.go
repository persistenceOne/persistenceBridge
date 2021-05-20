package contracts

import (
	"log"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/persistenceOne/persistenceCore/kafka/utils"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
)

var TokenWrapper = Contract{
	name:    "TOKEN_WRAPPER",
	address: constants.TokenWrapperAddress,
	abi:     abi.ABI{},
	methods: map[string]func(kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec, arguments []interface{}) error{
		constants.TokenWrapperWithdrawUTokens: onWithdrawUTokens,
	},
}

func onWithdrawUTokens(kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec, arguments []interface{}) error {
	// ercAddress := arguments[0].(common.Address)
	amount := arguments[1].(*big.Int)
	atomAddress, err := sdkTypes.AccAddressFromBech32(arguments[2].(string))
	if err != nil {
		return err
	}
	sendCoinMsg := bankTypes.NewMsgSend(constants.PSTakeAddress, atomAddress, sdkTypes.NewCoins(sdkTypes.NewCoin(constants.PSTakeDenom, sdkTypes.NewInt(amount.Int64()))))
	msgBytes, err := protoCodec.MarshalInterface(sdkTypes.Msg(sendCoinMsg))
	err = utils.ProducerDeliverMessage(msgBytes, utils.ToTendermint, kafkaState.Producer)
	if err != nil {
		log.Printf("Failed to add msg to kafka queue: %s\n", err.Error())
		return err
	}
	return nil
}
