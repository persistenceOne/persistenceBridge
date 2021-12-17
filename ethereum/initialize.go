/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package ethereum

import (
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	contracts2 "github.com/persistenceOne/persistenceBridge/ethereum/contracts"
)

func init() {
	contracts2.LiquidStaking.SetABI(constants2.LiquidStakingABI)
	contracts2.TokenWrapper.SetABI(constants2.TokenWrapperABI)
}
