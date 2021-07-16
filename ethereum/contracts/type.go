package contracts

import (
	"encoding/hex"
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type ContractI interface {
	GetName() string
	GetAddress() common.Address
	GetABI() abi.ABI
	SetABI(contractABIString string)
	GetMethods() map[string]func(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec, arguments []interface{}) error
	GetMethodAndArguments(inputData []byte) (*abi.Method, []interface{}, error)
}

type Contract struct {
	name    string
	address common.Address
	abi     abi.ABI
	methods map[string]func(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec, arguments []interface{}) error
}

var _ ContractI = &Contract{}

func (contract *Contract) GetName() string {
	return contract.name
}

func (contract *Contract) GetAddress() common.Address {
	return contract.address
}

func (contract *Contract) GetABI() abi.ABI {
	return contract.abi
}

func (contract *Contract) GetMethods() map[string]func(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec, arguments []interface{}) error {
	return contract.methods
}

func (contract *Contract) SetABI(contractABIString string) {
	contractABI, err := abi.JSON(strings.NewReader(contractABIString))
	if err != nil {
		log.Fatalln("Unable to decode abi:  " + err.Error())
	}
	contract.abi = contractABI
}

func (contract *Contract) GetMethodAndArguments(inputData []byte) (*abi.Method, []interface{}, error) {
	txData := hex.EncodeToString(inputData)
	if txData[:2] == "0x" {
		txData = txData[2:]
	}

	decodedSig, err := hex.DecodeString(txData[:8])
	if err != nil {
		logging.Fatal("Unable decode method ID (decodeSig) of", contract.name, "Error:", err)
	}

	method, err := contract.abi.MethodById(decodedSig)
	if err != nil {
		logging.Fatal("Unable to fetch method of", contract.name, "Error:", err)
	}

	decodedData, err := hex.DecodeString(txData[8:])
	if err != nil {
		logging.Fatal("Unable to decode input data of", contract.name, "Error:", err)
	}

	arguments, err := method.Inputs.Unpack(decodedData)
	return method, arguments, err
}
