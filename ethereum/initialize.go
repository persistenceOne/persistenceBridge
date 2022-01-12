/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package ethereum

import (
	"github.com/persistenceOne/persistenceBridge/ethereum/abi/liquidStaking"
	"github.com/persistenceOne/persistenceBridge/ethereum/abi/tokenWrapper"
	"github.com/persistenceOne/persistenceBridge/ethereum/contracts"
)

func init() {
	contracts.LiquidStaking.SetABI(liquidStaking.LiquidStakingMetaData.ABI)
	contracts.TokenWrapper.SetABI(tokenWrapper.TokenWrapperMetaData.ABI)
}
