package contracts

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/accounts/abi"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
)

var STokens = Contract{
	name:    "S_TOKENS",
	address: constants2.STokenAddress,
	abi:     abi.ABI{},
	methods: map[string]func(kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec, arguments []interface{}) error{},
}
