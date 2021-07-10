package contracts

import (
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
)

var STokens = Contract{
	name:    "S_TOKENS",
	address: common.HexToAddress(constants2.STokenAddress),
	abi:     abi.ABI{},
	methods: map[string]func(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec, arguments []interface{}) error{},
}
