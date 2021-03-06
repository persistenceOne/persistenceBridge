/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package contracts

import (
	"math/big"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

var TokenWrapper = Contract{
	name:    "TOKEN_WRAPPER",
	address: common.HexToAddress(configuration.GetAppConfig().Ethereum.TokenWrapperAddress),
	abi:     abi.ABI{},
	methods: map[string]func(arguments []interface{}) (sdkTypes.Msg, common.Address, error){
		constants.TokenWrapperWithdrawUTokens: onWithdrawUTokens,
	},
}

func onWithdrawUTokens(arguments []interface{}) (sdkTypes.Msg, common.Address, error) {
	ercAddress := arguments[0].(common.Address)
	amount := sdkTypes.NewIntFromBigInt(arguments[1].(*big.Int))
	atomAddress, err := sdkTypes.AccAddressFromBech32(arguments[2].(string))
	if err != nil {
		return nil, common.Address{}, err
	}
	sendCoinMsg := &bankTypes.MsgSend{
		FromAddress: configuration.GetAppConfig().Tendermint.GetWrapAddress(),
		ToAddress:   atomAddress.String(),
		Amount:      sdkTypes.NewCoins(sdkTypes.NewCoin(configuration.GetAppConfig().Tendermint.Denom, amount)),
	}
	logging.Info("Received ETH Unwrap Tx from:", ercAddress.String(), "amount:", amount.String(), "toAddress:", atomAddress.String())
	return sendCoinMsg, ercAddress, nil
}
