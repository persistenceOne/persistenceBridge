/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package contracts

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

type ContractI interface {
	GetName() string
	GetAddress() common.Address
	SetAddress(common.Address)
	GetABI() abi.ABI
	SetABI(contractABIString string)
	GetSDKMsgAndSender() map[string]func(arguments []interface{}) (sdk.Msg, common.Address, error)
	GetMethodAndArguments(inputData []byte) (*abi.Method, []interface{}, error)
}

type Contract struct {
	name    string
	address common.Address
	abi     abi.ABI
	methods map[string]func(arguments []interface{}) (sdk.Msg, common.Address, error)
}

var _ ContractI = &Contract{}

func (contract *Contract) GetName() string {
	return contract.name
}

func (contract *Contract) GetAddress() common.Address {
	return contract.address
}

func (contract *Contract) SetAddress(address common.Address) {
	contract.address = address
}

func (contract *Contract) GetABI() abi.ABI {
	return contract.abi
}

func (contract *Contract) GetSDKMsgAndSender() map[string]func(arguments []interface{}) (sdk.Msg, common.Address, error) {
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
	if len(txData) > 9 {

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
	return nil, nil, fmt.Errorf("invalid input data")
}
