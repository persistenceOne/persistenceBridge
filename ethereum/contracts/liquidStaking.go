/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package contracts

import (
	"math/big"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

var LiquidStaking = Contract{
	name:    "LIQUID_STAKING",
	address: common.HexToAddress(constants2.LiquidStakingAddress),
	abi:     abi.ABI{},
	methods: map[string]func(arguments []interface{}) (sdkTypes.Msg, common.Address, error){
		constants2.LiquidStakingStake:   onStake,
		constants2.LiquidStakingUnStake: onUnStake,
	},
}

// nolint implements common interface
// nolint:unparam
func onStake(arguments []interface{}) (sdkTypes.Msg, common.Address, error) {
	ercAddress := arguments[0].(common.Address)
	amount := sdkTypes.NewIntFromBigInt(arguments[1].(*big.Int))

	stakeMsg := &stakingTypes.MsgDelegate{
		DelegatorAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
		ValidatorAddress: "",
		Amount:           sdkTypes.NewCoin(configuration.GetAppConfig().Tendermint.PStakeDenom, amount),
	}

	logging.Info("Received ETH Stake Tx from:", ercAddress.String(), "amount:", amount.String())

	return stakeMsg, ercAddress, nil
}

// nolint implements a common interface
// nolint:unparam
func onUnStake(arguments []interface{}) (sdkTypes.Msg, common.Address, error) {
	ercAddress := arguments[0].(common.Address)
	amount := sdkTypes.NewIntFromBigInt(arguments[1].(*big.Int))

	unStakeMsg := &stakingTypes.MsgUndelegate{
		DelegatorAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
		ValidatorAddress: "",
		Amount:           sdkTypes.NewCoin(configuration.GetAppConfig().Tendermint.PStakeDenom, amount),
	}

	logging.Info("Received ETH UnStake Tx from:", ercAddress.String(), "amount:", amount.String())

	return unStakeMsg, ercAddress, nil
}
